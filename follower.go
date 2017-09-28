package main

import (
	"time"
	log "github.com/sirupsen/logrus"
)

type Follower struct {
	stop_chan chan bool
	wait_append_request_timeout_index int
	wait_append_request_timeout chan int
}

func (follower *Follower) Start(state *State,
	be_c chan bool, be_f chan bool, be_l chan bool, packet_chan chan Packet)  {

	log.Info("Follower Start")

	follower.stop_chan = make(chan bool, 10)
	follower.wait_append_request_timeout_index = 0
	follower.wait_append_request_timeout = make(chan int, 100)

	go func() {

		follower.StartWaitAppendRequestTimeout()

	outer:
		for {
			select {
			case <- follower.wait_append_request_timeout:

				// --- VoteRequest ----
				state.User = "candidate"
				be_c <- true
			case req := <- packet_chan:
				log.Info("follower", "get", req)

				// --- AppendRequest ---
				if req.Name == "AppendRequest" {
					if req.Term > state.CurrentTerm {
						state.UpdateCurrentTerm(req.Term)
					} else if req.Term == state.CurrentTerm {
					} else {
					}
					follower.StartWaitAppendRequestTimeout()
					res := Packet{Name: "AppendResponse", Id: state.Ip, Term: state.CurrentTerm}
					SendPacket(req.Id, res)

				// --- VoteRequest ----
				} else if req.Name == "VoteRequest" {
					if req.Term > state.CurrentTerm {
						state.UpdateCurrentTerm(req.Term)
						state.Vote = 0
						res := Packet{Name: "VoteResponse", Id: state.Ip, Term: state.CurrentTerm, VoteGranted: true}
						SendPacket(req.Id, res)
					} else if req.Term == state.CurrentTerm {
						if state.Vote == 1 {
							state.Vote = 0
							res := Packet{Name: "VoteResponse", Id: state.Ip, Term: state.CurrentTerm, VoteGranted: true}
							SendPacket(req.Id, res)
						} else {
							res := Packet{Name: "VoteResponse", Id: state.Ip, Term: state.CurrentTerm, VoteGranted: false}
							SendPacket(req.Id, res)
						}
					} else {
						res := Packet{Name: "VoteResponse", Id: state.Ip, Term: state.CurrentTerm, VoteGranted: false}
						SendPacket(req.Id, res)
					}
				}
			case <- follower.stop_chan:
				break outer
			}
		}
	}()
}

func (follower *Follower) Stop()  {
	log.Info("Follower Stop")

	follower.stop_chan = make(chan bool, 10)
	follower.wait_append_request_timeout_index = 0
	follower.wait_append_request_timeout = make(chan int, 100)

	follower.stop_chan <- true
}

func (follower *Follower) StartWaitAppendRequestTimeout()  {
	go func() {
		follower.wait_append_request_timeout_index++
		index := follower.wait_append_request_timeout_index

		time.Sleep(1500 * time.Millisecond)

		follower.wait_append_request_timeout <- index
	}()
}