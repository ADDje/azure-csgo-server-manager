package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/kidoman/go-steam"
)

const REPLAY_CONTAINER_NAME = "replays"
const REPLAY_FOLDER = "/home/steam/steamcmd/cs_go/csgo/*.dem"

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

func ExportReplays(config Config, week string, vmName string, vmUsername string, vmPassword string) error {

	log.Printf("Exporting Replays for: %s with label: %s", vmName, week)

	vmInfo, err := GetVmProperties(config, vmName)
	if err != nil {
		return err
	}

	isOn := false
	for _, status := range *vmInfo.InstanceView.Statuses {
		if *status.Code == "PowerState/running" {
			isOn = true
			break
		}
	}

	if !isOn {
		log.Printf("VM not on")
		return errors.New("VM is not on")
	}

	ip, err := GetVmIpWithProperties(config, vmInfo)
	if err != nil {
		return nil
	}

	return exportReplaysViaSSH(config, week, *ip, vmUsername, vmPassword)
}

func exportReplaysViaSSH(config Config, week string, ip string, vmUsername string, vmPassword string) error {
	sshConfig := ssh.ClientConfig{
		User: vmUsername,
		Auth: []ssh.AuthMethod{
			ssh.Password(vmPassword),
		},
		Timeout: time.Second * 10,
	}

	ip = fmt.Sprintf("%s:%d", ip, 22)
	log.Printf("Connecting to: %s", ip)
	conn, err := ssh.Dial("tcp", ip, &sshConfig)
	if err != nil {
		log.Printf("Couldn't connect to server: %s", err)
		return errors.New("Couldn't connect to server")
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Printf("Couldn't create session: %s", err)
		return errors.New("Couldn't create session")
	}
	defer session.Close()

	session.Stdout = mw
	session.Stderr = mw

	modes := ssh.TerminalModes{}

	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		return errors.New("Request for terminal failed")
	}

	replayFolder := week + "/$(hostname)"
	script := "#!/bin/bash\n\n" +
		"export AZURE_STORAGE_ACCOUNT=" + config.AzureStorageServer + "\n" +
		"export AZURE_STORAGE_ACCESS_KEY=" + config.AzureStorageKey + "\n" +

		"azure telemetry --disable\n" +

		"for f in " + REPLAY_FOLDER + "\n" +
		"do\n" +
		"echo \"Uploading $f file... ($(hostname))\"\n" +
		"azure storage blob upload -q $f " + REPLAY_CONTAINER_NAME + " \"" + replayFolder + "/$(basename $f)\"\n" +
		"done"

	session.Run("echo '" + script + "' > ~/upload.sh && chmod +x ~/upload.sh && ~/upload.sh")

	log.Printf("Export for %s complete", ip)

	return nil
}
