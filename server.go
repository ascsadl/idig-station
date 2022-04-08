package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"idig-station/serial"
)

func SerialServer(device string, port int) {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("Failed to start serial server: %s", err)
		return
	}

	go func() {
		defer listener.Close()

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Print(err)
				continue
			}

			client := conn.RemoteAddr().String()
			log.Printf("Received connection from %s", client)

			options := serial.OpenOptions{
				PortName:        device,
				BaudRate:        9600,
				DataBits:        8,
				StopBits:        1,
				MinimumReadSize: 1,
			}
			tty, err := serial.Open(options)
			if err != nil {
				log.Printf("Failed to open serial port: %v", err)
				log.Printf("Disconnected %s", client)
				conn.Close()
				continue
			}

			rtty := bufio.NewReader(tty)
			wtty := bufio.NewWriter(tty)
			rnet := bufio.NewReader(conn)
			wnet := bufio.NewWriter(conn)

			go leicaToiDig(rtty, wnet)
			iDigToLeica(rnet, wtty)

			tty.Close()
			conn.Close()

			log.Printf("Disconnected %s", client)
		}
	}()
}

func leicaToiDig(rtty *bufio.Reader, wnet *bufio.Writer) {
	var err error
	for {
		var line string
		line, err = rtty.ReadString('\n')
		if err != nil {
			err = nil
			break
		}
		msg := strings.TrimRight(line, "\r\n")
		if msg == "?" {
			continue // skip
		}

		log.Printf("Leica -> iDig: %q", msg)

		_, err = wnet.WriteString(msg + "\r\n")
		if err != nil {
			break
		}
		err = wnet.Flush()
		if err != nil {
			break
		}
	}
	if err != nil {
		log.Printf("Leica -> iDig: %s", err)
	}
}

func iDigToLeica(rnet *bufio.Reader, wtty *bufio.Writer) {
	var err error
	for {
		var line string
		line, err = rnet.ReadString('\n')
		if err != nil {
			err = nil
			break
		}
		msg := strings.TrimRight(line, "\r\n")
		log.Printf("iDig -> Leica: %q", msg)

		_, err = wtty.WriteString(msg + "\r\n")
		if err != nil {
			break
		}
		err = wtty.Flush()
		if err != nil {
			break
		}
	}
	if err != nil {
		log.Printf("iDig -> Leica: %s", err)
	}
}
