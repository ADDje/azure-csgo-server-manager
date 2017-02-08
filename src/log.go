package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Number of messages to store, to provide log history to new clients.
	bufferSize = 1024
)

var upgrader = websocket.Upgrader{

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var logServer *http.ServeMux
var messageBuffer [][]byte

type LogWriter struct {
	hub *Hub
}

func (l *LogWriter) SetHub(h *Hub) {
	l.hub = h
}

func (l *LogWriter) Write(data []byte) (int, error) {
	l.hub.broadcast <- data
	messageBuffer = append(messageBuffer, data)
	return len(data), nil
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			client.send <- getCurrentLog()
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func getCurrentLog() []byte {
	var outBuffer []byte

	for _, msg := range messageBuffer {
		for _, b := range msg {
			outBuffer = append(outBuffer, b)
		}
	}

	return outBuffer
}

func handleLog(hub *Hub, w http.ResponseWriter, r *http.Request) {

	if err := Auth.aaa.Authorize(w, r, true); err != nil {
		log.Printf("Unauthenticated WS request %s", r.RemoteAddr)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WS Error: %s", err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	go client.writePump()
}

func SetupLogWs() io.Writer {
	hub := newHub()
	go hub.run()

	logWriter := new(LogWriter)
	logWriter.SetHub(hub)

	logServer = http.NewServeMux()
	logServer.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleLog(hub, w, r)
	})

	return logWriter
}

func RunLogWs(config Config) {
	if logServer == nil {
		log.Fatal("Setup must be called before running Log Websocket Server")
	}

	addr := fmt.Sprintf("%s:%d", config.ServerIP, config.WebsocketPort)
	if config.UseSsl {
		log.Printf("Starting log WS server on: https://%s", addr)
		log.Fatal(http.ListenAndServeTLS(addr, config.SslCert, config.SslKey, logServer))
	} else {
		log.Printf("Starting log WS server on: http://%s", addr)
		log.Fatal(http.ListenAndServe(addr, logServer))
	}
}
