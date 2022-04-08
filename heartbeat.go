package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"time"
)

const HeartBeatInterval = 2 * time.Second

func HeartBeatGenerator(name string, port int) {
	msg := makeHeartBeat(name, port)

	for {
		time.Sleep(HeartBeatInterval)

		for _, bcast := range getInterfaceBroadcastAddrs() {
			addr := bcast + ":55555"
			conn, err := net.Dial("udp", addr)
			if err != nil {
				log.Printf("Error sending heartbeat: %s: %s", addr, err)
				continue
			}

			if _, err := conn.Write(msg); err != nil {
				log.Printf("Error sending heartbeat: %s: %s", addr, err)
			}

			conn.Close()
		}
	}
}

func makeHeartBeat(name string, port int) []byte {
	// Heartbeat packet is 110 bytes:
	// 0-5		MAC address of AP that we are Associated with (for location)
	// 6		Channel we are on
	// 7		RSSI
	// 8-9		local TCP port# (for connecting into the WiSnap device)
	// 10-13	RTC value (MSB first to LSB last)
	// 14-15	Battery Voltage on Pin 20 in millivolts (2755 for example)
	// 16-17	Value of the GPIO pins
	// 18-31	ASCII time
	// 32-59	Version string with date code
	// 60-91	Programmable Device ID string (set option deviceid <string>)
	// 92-93	Boot time in milliseconds
	// 94-110	Voltage readings of Sensors 0 through 7 (enabled with set opt format <mask>)

	var buf bytes.Buffer

	writeZeros(&buf, 8)
	writeUInt16(&buf, uint16(port))
	writeZeros(&buf, 50)
	writeString(&buf, name, 32)
	writeZeros(&buf, 18)

	return buf.Bytes()
}

func writeZeros(buf *bytes.Buffer, n int) {
	buf.Write(make([]byte, n))
}

func writeUInt16(buf *bytes.Buffer, n uint16) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, n)
	buf.Write(b)
}

func writeString(buf *bytes.Buffer, s string, n int) {
	buf.WriteString(s)
	buf.Write(make([]byte, n-len(s)))
}

func getInterfaceBroadcastAddrs() []string {
	var bcastAddrs []string
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip, ipNet, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}

			if ip.To4() == nil || !ip.IsGlobalUnicast() {
				continue
			}

			bcast := ipNet.IP
			for i, b := range ipNet.Mask {
				bcast[i] |= ^b
			}

			bcastAddrs = append(bcastAddrs, bcast.String())
		}
	}
	return bcastAddrs
}
