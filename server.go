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
	log.Info("Server Start with " + ip)

	server.PacketChan = make(chan Packet, 100)

	go func() {
		l, err := net.Listen("tcp", ip)
		if err != nil {
			log.Error("server not start: ", err.Error())
		}

		for {
			conn, err := l.Accept()
			if err != nil {
				continue
			}
			server.HandleRequest(conn)
		}
	}()
}


func (server *Server) Stop ()  {
	log.Info("Server Stop")
}

func (server *Server) HandleRequest(conn net.Conn) {

	go func() {
		defer conn.Close()

		for {
			buf := bufio.NewReader(conn)
			data, err := buf.ReadString('\n')
			if err != nil {
				log.Error(err.Error())
				continue
			}

			var packet Packet
			err = json.Unmarshal([]byte(data), &packet)
			if err != nil {
				log.Error(err.Error())
				continue
			}

			server.PacketChan <- packet
		}
	}()
}

