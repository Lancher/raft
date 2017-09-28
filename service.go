package main

import (
	log "github.com/sirupsen/logrus"
	"flag"
	"strings"
)

type Service struct {}

func (service *Service) Start()  {
	log.Info("Service Start")

	ip := flag.String("ip", "not-a-valid-ip", "self ip")
	ips := flag.String("ips", "not-a-valid-ip", "ip1,ip2,ip3....")
	flag.Parse()

	state := State{
		Ip: *ip,
		Ips: strings.Split(*ips, ","),
		CurrentTerm: 0,
		User: "follower",
	}

	// TODO load persistent configuration

	server := Server{}
	server.Start(state.Ip)

	be_c := make(chan bool, 10)
	be_f := make(chan bool, 10)
	be_l := make(chan bool, 10)

	candidate := Candidate{}
	follower := Follower{}
	leader := Leader{}

	if state.User == "candidate"  {
		be_c <- true
	} else if state.User == "follower" {
		be_f <- true
	} else if state.User == "leader" {
		be_l <- true
	}

	for {
		select {
		case <- be_c:
			leader.Stop()
			follower.Stop()
			candidate.Start(&state, be_c, be_f, be_l, server.PacketChan)
		case <- be_f:
			leader.Stop()
			candidate.Stop()
			follower.Start(&state, be_c, be_f, be_l, server.PacketChan)
		case <- be_l:
			follower.Stop()
			candidate.Stop()
			leader.Start(&state, be_c, be_f, be_l, server.PacketChan)
		}
	}
}

func (service *Service) Stop()  {
	log.Info("Service Stop")
}
