package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Username            string `json:"username"`
	Password            string `json:"password"`
	LogFolder           string `json:"log_folder"`
	DatabaseFile        string `json:"database_file"`
	CookieEncryptionKey string `json:"cookie_encryption_key"`
	//MaxUploadSize       int64  `json:"max_upload_size"`
	ServerIP            string `json:"server_ip"`
	ServerPort          int    `json:"server_port"`
	UseSsl              bool   `json:"use_ssl"`
	SslCert             string `json:"ssl_cert"`
	SslKey              string `json:"ssl_key"`
	AzureClientID       string `json:"azure_client_id"`
	AzureClientSecret   string `json:"azure_client_secret"`
	AzureSubscriptionID string `json:"azure_subscription_id"`
	AzureTenantID       string `json:"azure_tenant_id"`
	ResourceGroup       string `json:"resource_group"`
	UseCloudStorage     bool   `json:"use_cloud_storage"`
	AzureStorageServer  string `json:"azure_storage_server"`
	AzureStorageKey     string `json:"azure_storage_key"`
	AzureSASToken       string `json:"azure_sas_token"`
	VMVhdStorageServer  string `json:"vm_vhd_storage_server"`
	VMVhdStorageKey     string `json:"vm_vhd_storage_key"`
	ExternalApiKey      string `json:"external_api_key"`
	WebsocketPort       int    `json:"websocket_port"`

	ConfFile     string `json:"-"`
	WebSocketKey string `json:"-"`
	// Used to force not SSL inside tomcat
	OverrideSSL bool `json:"-"`
}

var (
	config Config
	Auth   *AuthHTTP
	mw     io.Writer
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// Loads server configuration files
// JSON config file contains default values,
// config file will overwrite any provided flags
func loadServerConfig(f string) {
	// Preserve some configs
	oldPort := config.ServerPort
	oldIP := config.ServerIP

	file, err := os.Open(f)
	failOnError(err, "Error loading config file.")

	err = json.NewDecoder(file).Decode(&config)

	if oldPort != 0 {
		config.ServerPort = oldPort
	}
	if oldIP != "" {
		config.ServerIP = oldIP
	}
	if config.OverrideSSL {
		config.UseSsl = false
	}
}

func parseFlags() {
	confFile := flag.String("conf", "./conf.json", "Specify location of Azure CSGO Server Manager config file.")
	webserverIP := flag.String("host", "0.0.0.0", "Specify IP for webserver to listen on.")
	webserverPort := flag.Int("port", 8090, "Specify a port for the server.")
	//serverMaxUpload := flag.Int64("max-upload", 1024*1024*20, "Maximum filesize for uploaded files (default 20MB).")

	flag.Parse()

	config.ConfFile = *confFile
	config.ServerIP = *webserverIP
	config.ServerPort = *webserverPort
	//config.MaxUploadSize = *serverMaxUpload

	port := os.Getenv("HTTP_PLATFORM_PORT")
	log.Printf("HTTP_PLATFORM_PORT: %s", port)
	if port != "" {
		myPort, err := strconv.Atoi(port)
		if err == nil {
			config.ServerPort = myPort
			config.OverrideSSL = true
		} else {
			log.Printf("Could not read port from HTTP_PLATFORM_PORT %s", err)
		}
		config.UseSsl = false
	}

	log.Printf("Server port: %d", config.ServerPort)
}

func setupLogging() {

	// Setup WebSocket key
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	config.WebSocketKey = fmt.Sprintf("%X", b)

	logFolder := config.LogFolder
	if logFolder == "" {
		log.Printf("WARNING: Using local log folder as none specified")
		logFolder = "./log"
	}

	logWriter := SetupLogWs()
	go RunLogWs(config)

	mw = io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Dir:     logFolder,
		MaxSize: 1024 * 100,
		MaxAge:  7,
	}, logWriter)
	log.SetOutput(mw)
}

func main() {
	parseFlags()
	loadServerConfig(config.ConfFile)

	setupLogging()

	// Initialize authentication system
	Auth = initAuth()
	Auth.CreateAuth(config.DatabaseFile, config.CookieEncryptionKey)
	Auth.CreateOrUpdateUser(config.Username, config.Password, "admin", "")

	router := NewRouter()

	addr := fmt.Sprintf("%s:%d", config.ServerIP, config.ServerPort)
	if config.UseSsl {
		log.Printf("Starting server on: https://%s", addr)
		log.Fatal(http.ListenAndServeTLS(addr, config.SslCert, config.SslKey, router))
	} else {
		log.Printf("Starting server on: http://%s", addr)
		log.Fatal(http.ListenAndServe(addr, router))
	}

}
