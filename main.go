package main

import (
	"log"
	"net"

	"github.com/ardean/qnap-lcd-display-manager/display"
)

func main() {
	log.Println("Starting QNAP LCD Display Manager...")

	ipAddresses, err := ListIPAddresses()
	panicCheck(err)

	currentIPAddressIndex := 0
	currentIPAddress := ipAddresses[currentIPAddressIndex]

	lcdAddress := display.Find()
	defer func() {
		if lcdAddress != nil {
			panicCheck((*lcdAddress).Close())
		}
	}()

	PrintIPAddress := func(ipAddress string) {
		if lcdAddress != nil {
			(*lcdAddress).Write(0, "IP:")
			(*lcdAddress).Write(1, ipAddress)
		} else {
			log.Println("IP:")
			log.Println(ipAddress)
		}
	}

	GetNextIPAddress := func(direction int) int {
		var nextIndex int

		ipAddressCount := len(ipAddresses)
		if ipAddressCount > 0 {
			nextIPAddressIndex := currentIPAddressIndex + direction
			if nextIPAddressIndex < 0 {
				nextIPAddressIndex = ipAddressCount - 1
			} else if nextIPAddressIndex >= ipAddressCount {
				nextIPAddressIndex = 0
			}
			nextIndex = nextIPAddressIndex
		} else {
			nextIndex = -1
		}

		return nextIndex
	}

	MoveIPAddress := func(direction int) {
		currentIPAddressIndex = GetNextIPAddress(direction)

		if currentIPAddressIndex != -1 {
			currentIPAddress = ipAddresses[currentIPAddressIndex]
			PrintIPAddress(currentIPAddress)
		}
	}

	if lcdAddress != nil {
		log.Println("Done")

		MoveIPAddress(0)

		(*lcdAddress).Listen(func(buttonIndex int, released bool) bool {
			if released {
				if buttonIndex == 1 {
					MoveIPAddress(1)
				} else if buttonIndex == 2 {
					MoveIPAddress(-1)
				}
			}

			return true
		})
	} else {
		log.Println("No display found! Exiting...")
	}
}

func panicCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func ListIPAddresses() ([]string, error) {
	ipAddresses := []string{}

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			var ip string
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP.String()
			case *net.IPAddr:
				ip = v.IP.String()
			}

			if len(ip) > 0 {
				ipAddresses = append(ipAddresses, ip)
			}
		}
	}

	return ipAddresses, nil
}
