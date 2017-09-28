package main

 import (
 	log "github.com/sirupsen/logrus"
 	"math/rand"
 	"time"
 )

type Candidate struct {
	stop_chan chan bool
	wait_vote_request_timeout_index int
	wait_vote_request_timeout chan int
}

func (candidate *Candidate) Start(state *State,
	be_c chan bool, be_f chan bool, be_l chan bool, packet_chan chan Packet)  {

	log.Info("Candidate Start")

	candidate.stop_chan = make(chan bool, 10)
	candidate.wait_vote_request_timeout_index = 0
	candidate.wait_vote_request_timeout = make(chan int, 100)

	go func() {

		candidate.SendVoteRequest(state)

	outer:
		for {
			select {
			case <- candidate.wait_vote_request_timeout:

				// --- VoteRequest ----
				candidate.SendVoteRequest(state)
			case req := <- packet_chan:
				log.Info("candidate", "get", req)

				// --- AppendRequest ---
				if req.Name == "AppendRequest" {
					if req.Term > state.CurrentTerm {
						state.UpdateCurrentTerm(req.Term)
						state.User = "follower"
						res := Packet{Name: "AppendResponse", Id: state.Ip, Term: state.CurrentTerm}
						SendPacket(req.Id, res)
						be_f <- true
					} else if req.Term == state.CurrentTerm {
						state.UpdateCurrentTerm(req.Term)
						state.User = "follower"
						res := Packet{Name: "AppendResponse", Id: state.Ip, Term: state.CurrentTerm}
						SendPacket(req.Id, res)
						be_f <- true
					} else {
						res := Packet{Name: "AppendResponse", Id: state.Ip, Term: state.CurrentTerm}
						SendPacket(req.Id, res)
					}

				// --- VoteRequest ---
				} else if req.Name == "VoteRequest" {
					if req.Term > state.CurrentTerm {
						state.UpdateCurrentTerm(req.Term)
						state.Vote = 0
						state.User = "follower"
						res := Packet{Name: "VoteResponse", Id: state.Ip, Term: state.CurrentTerm, VoteGranted: true}
						SendPacket(req.Id, res)
						be_f <- true
					} else if req.Term == state.CurrentTerm {
						res := Packet{Name: "VoteResponse", Id: state.Ip, Term: state.CurrentTerm, VoteGranted: false}
						SendPacket(req.Id, res)
					} else {
						res := Packet{Name: "VoteResponse", Id: state.Ip, Term: state.CurrentTerm, VoteGranted: false}
						SendPacket(req.Id, res)
					}

				// --- VoteResponse
				} else if req.Name == "VoteResponse" {
					if req.Term > state.CurrentTerm {
						state.UpdateCurrentTerm(req.Term)
						state.User = "follower"
						be_f <- true
					} else if req.Term == state.CurrentTerm {
						if req.VoteGranted == true { state.Vote ++}
					} else {
						if req.VoteGranted == true { state.Vote ++}
					}
					if state.Vote > len(state.Ips) / 2 {
						state.User = "leader"
						be_l <- true
					}
				}
			case <- candidate.stop_chan:
				break outer
			}
		}
	}()
}

func (candidate *Candidate) Stop()  {
	log.Info("Candidate Stop")
}

func (candidate *Candidate) StartWaitVoteRequestTimeout()  {
	go func() {
		candidate.wait_vote_request_timeout_index++
		index := candidate.wait_vote_request_timeout_index

		time.Sleep(time.Duration(rand.Intn(1000) + 1500) * time.Millisecond)

		candidate.wait_vote_request_timeout <- index
	}()
}

func (candidate *Candidate) SendVoteRequest(state *State)  {
	state.UpdateCurrentTerm(state.CurrentTerm + 1)
	go func() {
		for _, ip := range state.Ips {
			req := Packet{Name: "VoteRequest", Id: state.Ip, Term: state.CurrentTerm}
			SendPacket(ip, req)
		}
	}()
}