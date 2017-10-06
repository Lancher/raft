package main

import (
	log "github.com/sirupsen/logrus"
	"net"
	"bufio"
	"encoding/json"
)

type Server struct {
	PacketChan chan Packet
}


func (server *Server) Start (ip string)  {
	log.Info(ip, ": ", "server started")

	server.PacketChan = make(chan Packet, 100)

	go func() {
		l, err := net.Listen("tcp", ip)
		if err != nil {
			log.Error(ip, ": ", "server failed to start because ", err.Error())
		}

		for {
			conn, err := l.Accept()
			if err != nil {
				//log.Error(ip, ": ", "server failed to accept connections because ", err.Error())
				continue
			}
			server.HandleRequest(ip, conn)
		}
	}()
}


func (server *Server) HandleRequest(ip string, conn net.Conn) {

	go func() {
		defer conn.Close()

		for {
			buf := bufio.NewReader(conn)
			data, err := buf.ReadString('\n')
			if err != nil {
				//log.Error(ip, ": ", "server failed to accept connections because ", err.Error())
				return
			}

			var packet Packet
			err = json.Unmarshal([]byte(data), &packet)
			if err != nil {
				log.Error(ip, ": ", "server failed to decode packet because ", err.Error())
				return
			}
			server.PacketChan <- packet
		}
	}()
}

