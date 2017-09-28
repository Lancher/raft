package main

import (
	log "github.com/sirupsen/logrus"
)

type Leader struct {
	stop_chan chan bool
	send_append_request_timeout_index int
	send_append_request_timeout chan int
}

func (leader *Leader) Start(state *State,
	be_c chan bool, be_f chan bool, be_l chan bool, packet_chan chan Packet)  {
	log.Info("Leader Start")

	leader.stop_chan = make(chan bool, 10)
	leader.send_append_request_timeout_index = 0
	leader.send_append_request_timeout = make(chan int, 100)

	go func() {

		leader.SendAppendRequest(state)

	outer:
		for {
			select {
			case <- leader.send_append_request_timeout:
				// --- AppendRequest ---
				leader.SendAppendRequest(state)
			case req := <- packet_chan:
				log.Info("leader", "get", req)

				// --- AppendRequest ---
				if req.Name == "AppendRequest" {
					if req.Term > state.CurrentTerm {
						state.UpdateCurrentTerm(req.Term)
						res := Packet{Name: "AppendResponse", Id: state.Ip, Term: state.CurrentTerm}
						SendPacket(req.Id, res)
						state.User = "follower"
						be_f <- true
					} else if req.Term == state.CurrentTerm {
						res := Packet{Name: "AppendResponse", Id: state.Ip, Term: state.CurrentTerm}
						SendPacket(req.Id, res)
					} else {
						res := Packet{Name: "AppendResponse", Id: state.Ip, Term: state.CurrentTerm}
						SendPacket(req.Id, res)
					}

				// --- VoteRequest ---
				} else if req.Name == "VoteRequest" {
					if req.Term > state.CurrentTerm {
						state.UpdateCurrentTerm(req.Term)
						state.Vote = 0
						res := Packet{Name: "VoteResponse", Id: state.Ip, Term: state.CurrentTerm, VoteGranted: true}
						SendPacket(req.Id, res)
						state.User = "follower"
						be_f <- true
					} else if req.Term == state.CurrentTerm {
						res := Packet{Name: "VoteResponse", Id: state.Ip, Term: state.CurrentTerm, VoteGranted: false}
						SendPacket(req.Id, res)
					} else {
						res := Packet{Name: "VoteResponse", Id: state.Ip, Term: state.CurrentTerm, VoteGranted: false}
						SendPacket(req.Id, res)
					}

				// --- AppendResponse ---
				} else if req.Name == "AppendResponse" {
					if req.Term > state.CurrentTerm {
						state.UpdateCurrentTerm(req.Term)
						state.User = "follower"
						be_f <- true
					} else if req.Term == state.CurrentTerm {
					} else {
					}

				// --- VoteResponse ---
				} else if req.Name == "VoteResponse" {
					if req.Term > state.CurrentTerm {
						state.UpdateCurrentTerm(req.Term)
						state.User = "follower"
						be_f <- true
					} else if req.Term == state.CurrentTerm {
					} else {
					}
				}

			case <- leader.stop_chan:
				break outer
			}
		}
	}()
}

func (leader *Leader) Stop()  {
	log.Info("Leader Stop")

	leader.stop_chan = make(chan bool, 10)
	leader.send_append_request_timeout_index = 0
	leader.send_append_request_timeout = make(chan int, 100)

	leader.stop_chan <- true
}

func (leader *Leader) SendAppendRequest(state *State)  {
	go func() {
		for _, ip := range state.Ips {
			req := Packet{Name: "AppendRequest", Id: state.Ip, Term: state.CurrentTerm}
			SendPacket(ip, req)
		}
	}()
}
