package main

import (
	"log"

	"github.com/kidoman/go-steam"
)

func GetInfoForServer(serverAddress string) (*steam.InfoResponse, error) {

	server, err := steam.Connect(serverAddress)

	if err != nil {
		log.Printf("Error in steam getInfoForServer connect: %s", err)
		return nil, err
	}

	info, err := server.Info()

	if err != nil {
		log.Printf("Error in steam getInfoForServer info: %s", err)
		return nil, err
	}

	return info, nil
}
