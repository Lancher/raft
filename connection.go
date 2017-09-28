package main

import (
	"net"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

func SendPacket(ip string, packet Packet)  {
	go func() {
		conn, err := net.Dial("tcp", ip)
		if err != nil {
			log.Error(err.Error())
			return
		}
		defer conn.Close()

		bytes, err := json.Marshal(packet)
		newline := "\n"
		conn.Write(append(bytes, newline...))
	}()
}