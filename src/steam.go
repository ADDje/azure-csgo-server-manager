package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

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

func ExportReplays(config Config, vmName string, vmUsername string, vmPassword string) error {

	log.Printf("Exporting Replays for: %s", vmName)

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

	nic := *(*vmInfo.NetworkProfile.NetworkInterfaces)[0].ID
	nicParts := strings.Split(nic, "/")
	nicName := nicParts[len(nicParts)-1]

	nicDetails, err := GetNicDetails(config, nicName)
	if err != nil {
		return err
	}
	nicDetails2 := *nicDetails

	ip := (*nicDetails2.IPConfigurations)[0]
	pubIP := *ip.PublicIPAddress

	ipParts := strings.Split(*pubIP.ID, "/")
	ipID := ipParts[len(ipParts)-1]

	ipDetails, err := GetIpDetails(config, ipID)
	if err != nil {
		return err
	}

	return exportReplaysViaSSH(*ipDetails.IPAddress, vmUsername, vmPassword)
}

func exportReplaysViaSSH(ip string, vmUsername string, vmPassword string) error {
	config := ssh.ClientConfig{
		User: vmUsername,
		Auth: []ssh.AuthMethod{
			ssh.Password(vmPassword),
		},
		Timeout: time.Second * 10,
	}

	ip = fmt.Sprintf("%s:%d", ip, 22)
	log.Printf("Connecting to: %s", ip)
	conn, err := ssh.Dial("tcp", ip, &config)
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

	session.Run("cd /home/steam/ && ./upload.sh")

	log.Printf("Export for %s complete", ip)

	return nil
}
