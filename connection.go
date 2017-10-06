package main

import (
	"net"
	"encoding/json"
	//log "github.com/sirupsen/logrus"
)

func SendPacket(packet Packet)  {
	go func() {
		conn, err := net.Dial("tcp", packet.DST_ID)
		if err != nil {
			//log.Error(packet.Id, ": ", "client failed to dail to ", packet.DST_ID, " because ", err.Error())
			return
		}
		defer conn.Close()

		bytes, err := json.Marshal(packet)
		newline := "\n"
		conn.Write(append(bytes, newline...))
		//log.Error(packet.Id, ": ", "client dailed to ", packet.DST_ID)
	}()
}