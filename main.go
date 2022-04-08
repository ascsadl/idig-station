package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

func main() {
	nameFlag := flag.String("n", "", "Set iDig Station's name")
	devFlag := flag.String("d", "/dev/ttyUSB0", "Serial device")
	portFlag := flag.Int("p", 4000, "Port to listen on")
	flag.Parse()

	name := *nameFlag
	if name == "" {
		data, err := os.ReadFile("/etc/idig-station")
		if err == nil {
			name = strings.TrimSpace(string(data))
		}
	}
	if name == "" {
		name = "iDig Station"
	}

	SerialServer(*devFlag, *portFlag)
	log.Printf("Starting %s on port %d", name, *portFlag)

	HeartBeatGenerator(name, *portFlag)
}
