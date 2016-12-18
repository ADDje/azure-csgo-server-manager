package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Config struct {
	AzureClientID       string `json:"azure_client_id"`
	AzureClientSecret   string `json:"azure_client_secret"`
	AzureSubscriptionID string `json:"azure_subscription_id"`
	AzureTenantID       string `json:"azure_tenant_id"`
	ResourceGroup       string `json:"resource_group"`
	ServerIP            string `json:"server_ip"`
	ServerPort          string `json:"server_port"`
	MaxUploadSize       int64  `json:"max_upload_size"`
	Username            string `json:"username"`
	Password            string `json:"password"`
	DatabaseFile        string `json:"database_file"`
	CookieEncryptionKey string `json:"cookie_encryption_key"`
	UseCloudStorage     bool   `json:"use_cloud_storage"`
	AzureStorageServer  string `json:"azure_storage_server"`
	AzureStorageKey     string `json:"azure_storage_key"`

	ConfFile string `json:"-"`
}

var (
	config Config
	Auth   *AuthHTTP
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
	file, err := os.Open(f)
	failOnError(err, "Error loading config file.")

	err = json.NewDecoder(file).Decode(&config)
}

func parseFlags() {
	confFile := flag.String("conf", "./conf.json", "Specify location of Azure CSGO Server Manager config file.")
	webserverIP := flag.String("host", "0.0.0.0", "Specify IP for webserver to listen on.")
	webserverPort := flag.String("port", "8090", "Specify a port for the server.")
	serverMaxUpload := flag.Int64("max-upload", 1024*1024*20, "Maximum filesize for uploaded files (default 20MB).")

	flag.Parse()

	config.ConfFile = *confFile
	config.ServerIP = *webserverIP
	config.ServerPort = *webserverPort
	config.MaxUploadSize = *serverMaxUpload
}

func main() {
	parseFlags()
	loadServerConfig(config.ConfFile)

	// Initialize authentication system
	Auth = initAuth()
	Auth.CreateAuth(config.DatabaseFile, config.CookieEncryptionKey)
	Auth.CreateOrUpdateUser(config.Username, config.Password, "admin", "")

	router := NewRouter()

	fmt.Printf("Starting server on: %s:%s", config.ServerIP, config.ServerPort)
	log.Fatal(http.ListenAndServeTLS(config.ServerIP+":"+config.ServerPort, "certs/cert.pem", "certs/key.pem", router))
}
