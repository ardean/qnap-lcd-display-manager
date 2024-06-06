package main

import (
	"log"
	"net"
	"time"

	"github.com/ardean/qnap-lcd-display-manager/display"
)

func main() {
	log.Println("Starting QNAP LCD Display Manager...")

	lcdAddress := display.Find()
	if lcdAddress == nil {
		log.Println("No display found! Exiting...")
		return
	}

	var ipAddresses []string
	var currentIPAddressIndex int = -1
	var currentIPAddress string
	var lastButtonPressed time.Time
	var nextStandby time.Time
	var isInStandby bool = true

	lcd := *lcdAddress

	defer func() {
		if lcdAddress != nil {
			panicCheck(lcd.Close())
		}
	}()

	PrintIPAddress := func(ipAddress string) {
		if lcdAddress != nil {
			lcd.Write(0, "IP:")
			lcd.Write(1, ipAddress)
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
		} else {
			PrintIPAddress("-")
		}
	}

	RefreshIPAddresses := func() {
		var err error
		ipAddresses, err = ListIPAddresses()
		panicCheck(err)

		MoveIPAddress(0)
	}

	Standby := func() {
		if isInStandby {
			return
		}
		isInStandby = true
		lcd.Enable(false)
	}

	startStandbyWatcher := func() {
		go (func() {
			for {
				time.Sleep(1 * time.Second)
				if time.Now().After(nextStandby) {
					Standby()
					break
				}
			}
		})()
	}

	startIPAddressRefresher := func() {
		go (func() {
			for {
				RefreshIPAddresses()
				time.Sleep(1 * time.Minute)
				if time.Now().After(nextStandby) {
					break
				}
			}
		})()
	}

	Wakeup := func() {
		if !isInStandby {
			return
		}
		isInStandby = false

		lcd.Enable(true)
		startStandbyWatcher()
		startIPAddressRefresher()
	}

	OnButtonPress := func(buttonIndex int, released bool) bool {
		if released {
			wasInStandby := isInStandby

			Wakeup()

			if wasInStandby {
				return true
			}

			if buttonIndex == 1 {
				MoveIPAddress(1)
			} else if buttonIndex == 2 {
				MoveIPAddress(-1)
			} else {
				Standby()
			}
		}

		lastButtonPressed = time.Now()
		nextStandby = lastButtonPressed.Add(10 * time.Second)

		return true
	}

	Wakeup()

	log.Println("Done")

	lcd.Listen(OnButtonPress)
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
	for _, iface := range ifaces {
		if (iface.Flags & net.FlagRunning) != net.FlagRunning {
			continue
		}

		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			var ip string
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.IsLoopback() {
					continue
				}
				ip = v.IP.String()
			case *net.IPAddr:
				if v.IP.IsLoopback() {
					continue
				}
				ip = v.IP.String()
			}

			if len(ip) > 0 {
				ipAddresses = append(ipAddresses, ip)
			}
		}
	}

	return ipAddresses, nil
}
