package main

 import (
 	log "github.com/sirupsen/logrus"
 	"math/rand"
 	"time"
 )

func CandidateHandleReceivePacket(state *State, rec Packet)  {

	// AppendRequest
	if rec.Name == "AppendRequest" {
		if state.CurrentTerm < rec.Term {
			// log
			log.Info(state.Ip, ": ", "candidate got AppendRequest with greater term and then become follower (normal)")

			// user, term, vote
			state.User = "follower"
			state.CurrentTerm = rec.Term
			state.Vote = 1

			// send packet
			packet := Packet{Name: "AppendResponse", Id: state.Ip, DST_ID: rec.Id, Term: state.CurrentTerm}
			state.SendPacketChan <- packet

			// become follower event
			state.EventChan <- "become_follower"
		} else if state.CurrentTerm == rec.Term {
			// log
			log.Info(state.Ip, ": ", "candidate got AppendRequest with the same term and then become follower (normal)")

			// user, term, vote
			state.User = "follower"
			state.CurrentTerm = rec.Term
			state.Vote = 1

			// send packet
			packet := Packet{Name: "AppendResponse", Id: state.Ip, DST_ID: rec.Id, Term: state.CurrentTerm}
			state.SendPacketChan <- packet

			// become follower event
			state.EventChan <- "become_follower"
		} else if state.CurrentTerm > rec.Term {
			// log
			log.Info(state.Ip, ": ", "candidate got AppendRequest with smaller term (normal)")

			// send packet
			packet := Packet{Name: "AppendResponse", Id: state.Ip, DST_ID: rec.Id, Term: state.CurrentTerm}
			state.SendPacketChan <- packet
		}

	// AppendResponse
	} else if rec.Name == "AppendResponse" {
		if state.CurrentTerm < rec.Term {
			// log
			log.Warning(state.Ip, ": ", "candidate got AppendResponse with greater term (should not happened)")
		} else if state.CurrentTerm == rec.Term {
			// log
			log.Warning(state.Ip, ": ", "candidate got AppendResponse with the same term (should not happened)")
		} else if state.CurrentTerm > rec.Term {
			// log
			log.Info(state.Ip, ": ", "candidate got AppendResponse with smaller term (ignore late response)")
		}

	// VoteRequest
	}else if rec.Name == "VoteRequest" {
		if state.CurrentTerm < rec.Term {
			// log
			log.Info(state.Ip, ": ", "candidate got VoteRequest with greater term and then become follower (normal)")

			// user, term, vote
			state.User = "follower"
			state.CurrentTerm = rec.Term
			state.Vote = 0

			// send packet
			log.Info("VoteResponse ", state.Ip, "->", rec.Id)
			packet := Packet{Name: "VoteResponse", Id: state.Ip, DST_ID: rec.Id, Term: state.CurrentTerm, VoteGranted: true}
			state.SendPacketChan <- packet

			// become follower event
			state.EventChan <- "become_follower"
		} else if state.CurrentTerm == rec.Term {
			// log
			log.Info(state.Ip, ": ", "candidate got VoteRequest with the same term (normal)")

			// send packet
			packet := Packet{Name: "VoteResponse", Id: state.Ip, DST_ID: rec.Id, Term: state.CurrentTerm, VoteGranted: false}
			state.SendPacketChan <- packet
		} else if state.CurrentTerm > rec.Term {
			// log
			log.Info(state.Ip, ": ", "candidate got VoteRequest with smaller term (normal)")

			// send packet
			packet := Packet{Name: "VoteResponse", Id: state.Ip, DST_ID: rec.Id, Term: state.CurrentTerm, VoteGranted: false}
			state.SendPacketChan <- packet
		}

	// VoteResponse
	}else if rec.Name == "VoteResponse" {
		if state.CurrentTerm < rec.Term {
			// log
			log.Info(state.Ip, ": ", "candidate got VoteResponse with greater term and then become follower (normal)")

			// user, term, vote
			state.User = "follower"
			state.CurrentTerm = rec.Term
			state.Vote = 1

			// become follower event
			state.EventChan <- "become_follower"
		} else if state.CurrentTerm == rec.Term {
			// log
			log.Info(state.Ip, ": ", "candidate got VoteResponse with the same term (normal)")

			if rec.VoteGranted == true { state.Vote ++ }
			if state.Vote > (len(state.Ips) + 1) / 2 {
				// user, term, vote
				state.User = "leader"
				state.CurrentTerm = rec.Term
				state.Vote = 1

				// become leader event
				state.EventChan <- "become_leader"
			}
		} else if state.CurrentTerm > rec.Term {
			// log
			log.Info(state.Ip, ": ", "candidate got VoteResponse with smaller term (ignore late response)")
		}
	}
}

func CandidateStartSendVoteRequestTimeout(state *State)  {

	state.send_vote_request_timeout_index ++
	index := state.send_vote_request_timeout_index

	go func() {
		if state.User != "candidate" { return }
		time.Sleep(time.Duration(rand.Intn(1000) + 1500) * time.Millisecond)

		if state.User != "candidate" { return }
		if state.send_vote_request_timeout_index != index { return }
		state.EventChan <- "send_vote_request_timeout"
	}()
}

func CandidateSendVoteRequest(state *State)  {
	if state.User != "candidate" { return }

	// log
	log.Info(state.Ip, ": ", "candidate sent VoteRequest to ", state.Ips)

	// user, term, vote
	state.User = "candidate"
	state.CurrentTerm ++
	state.Vote = 1

	go func() {
		for _, ip := range state.Ips {
			if state.User != "candidate" { return }
			packet := Packet{Name: "VoteRequest", Id: state.Ip, DST_ID:ip, Term: state.CurrentTerm}
			state.SendPacketChan <- packet
		}
	}()
}


