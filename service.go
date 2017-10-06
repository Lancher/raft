package main

import (
	log "github.com/sirupsen/logrus"
)




type State struct {
	User string
	CurrentTerm int
	Ip string
	Ips [] string
	Vote int

	EventChan chan string
	ReceivePacketChan chan Packet
	SendPacketChan chan Packet

	send_vote_request_timeout_index int   // candidate
	wait_append_request_timeout_index int // follower
	send_append_request_timeout_index int // leader
}

type Service struct {
}

func (service *Service) Start(ip string, ips [] string)  {

	// log setting
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)

	// tcp server
	server := Server{}
	server.Start(ip)

	// state
	state := State{
		User: "follower",
		CurrentTerm: 0,
		Ip: ip,
		Ips: ips,
		Vote: 1,

		EventChan: make(chan string, 100),
		ReceivePacketChan: server.PacketChan,
		SendPacketChan: make(chan Packet, 100),

		send_vote_request_timeout_index: 0,
		wait_append_request_timeout_index: 0,
		send_append_request_timeout_index: 0,
	}

	// TODO load persistent configurations in disk

	// start user
	if state.User == "candidate" {
		state.EventChan <- "become_candidate"
	} else if state.User == "follower" {
		state.EventChan <- "become_follower"
	} else if state.User == "leader" {
		state.EventChan <- "become_leader"
	}

	// event loop
	for {
		select {
		case event := <- state.EventChan:
			// candidate
			if event == "become_candidate" {
				log.Info(state.Ip, ": ", "start as ", state.User)
				CandidateStartSendVoteRequestTimeout(&state)
				CandidateSendVoteRequest(&state)
			// follower
			} else if event == "become_follower" {
				log.Info(state.Ip, ": ", "start as ", state.User)
				FollowerStartWaitAppendRequestTimeout(&state)
			// leader
			} else if event == "become_leader" {
				log.Info(state.Ip, ": ", "start as ", state.User)
				LeaderStartSendAppendRequestTimeout(&state)
				LeaderSendAppendRequest(&state)
			// candidate
			} else if event == "send_vote_request_timeout" {
				CandidateStartSendVoteRequestTimeout(&state)
				CandidateSendVoteRequest(&state)
			// follower
			} else if event == "wait_append_request_timeout" {
				FollowerBecomeCandidate(&state)
			// leader
			} else if event == "send_append_request_timeout" {
				LeaderStartSendAppendRequestTimeout(&state)
				LeaderSendAppendRequest(&state)
			}

		case rec := <- state.ReceivePacketChan:
			if state.User == "candidate" {
				CandidateHandleReceivePacket(&state, rec)
			} else if state.User == "follower" {
				FollowerHandleReceivePacket(&state, rec)
			} else if state.User == "leader" {
				LeaderHandleReceivePacket(&state, rec)
			}
		case snd := <- state.SendPacketChan:
			if snd.Name == "AppendRequest" {
				if state.User == "leader" {
					SendPacket(snd)
				}
			} else if snd.Name == "AppendResponse" {
				SendPacket(snd)
			} else if snd.Name == "VoteRequest" {
				if state.User == "candidate" {
					SendPacket(snd)
				}
			} else if snd.Name == "VoteResponse" {
				SendPacket(snd)
			}
		}
	}
}

func (service *Service) Stop()  {
}
