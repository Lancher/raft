package main

import (
	log "github.com/sirupsen/logrus"
	"time"
)


func LeaderHandleReceivePacket(state *State, rec Packet)  {
	// AppendRequest
	if rec.Name == "AppendRequest" {
		if state.CurrentTerm < rec.Term {
			// log
			log.Info(state.Ip, ": ", "leader got AppendRequest with the greater term and then become follower (normal)")

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
			log.Warning(state.Ip, ": ", "leader got AppendRequest with the same term (should not happened)")
		} else if state.CurrentTerm > rec.Term {
			//log
			log.Info(state.Ip, ": ", "leader got AppendRequest with the smaller term (normal)")

			// send packet
			packet := Packet{Name: "AppendResponse", Id: state.Ip, DST_ID: rec.Id, Term: state.CurrentTerm}
			state.SendPacketChan <- packet
		}

		// AppendResponse
	} else if rec.Name == "AppendResponse" {
		if state.CurrentTerm < rec.Term {
			//log
			log.Info(state.Ip, ": ", "leader got AppendResponse with the greater term and then become follower (normal)")

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
			log.Info(state.Ip, ": ", "leader got AppendResponse with the same term (normal)")
		} else if state.CurrentTerm > rec.Term {
			//log
			log.Info(state.Ip, ": ", "leader got AppendResponse with smaller term (ignore late response)")
		}

	// VoteRequest
	}else if rec.Name == "VoteRequest" {
		if state.CurrentTerm < rec.Term {
			//log
			log.Info(state.Ip, ": ", "leader got VoteRequest with the greater term and then become follower (normal)")

			// user, term, vote
			state.User = "follower"
			state.CurrentTerm = rec.Term
			state.Vote = 0

			// send packet
			packet := Packet{Name: "VoteResponse", Id: state.Ip, DST_ID: rec.Id, Term: state.CurrentTerm, VoteGranted: true}
			state.SendPacketChan <- packet

			// become follower event
			state.EventChan <- "become_follower"
		} else if state.CurrentTerm == rec.Term {
			//log
			log.Info(state.Ip, ": ", "leader got VoteRequest with the same term (normal)")

			// send packet
			packet := Packet{Name: "VoteResponse", Id: state.Ip, DST_ID: rec.Id, Term: state.CurrentTerm, VoteGranted: false}
			state.SendPacketChan <- packet
		} else if state.CurrentTerm > rec.Term {
			//log
			log.Info(state.Ip, ": ", "leader got VoteRequest with the smaller term (normal)")

			// send packet
			packet := Packet{Name: "VoteResponse", Id: state.Ip, DST_ID: rec.Id, Term: state.CurrentTerm, VoteGranted: false}
			state.SendPacketChan <- packet
		}

	// VoteResponse
	}else if rec.Name == "VoteResponse" {
		if state.CurrentTerm < rec.Term {
			log.Warning(state.Ip, ": ", "leader got VoteResponse with greater term (should not happened)")
		} else if state.CurrentTerm == rec.Term {
			log.Warning(state.Ip, ": ", "leader got VoteResponse with the same term (should not happened)")
		} else if state.CurrentTerm > rec.Term {
			log.Info(state.Ip, ": ", "leader got VoteResponse with smaller term (ignore late response)")
		}
	}
}

func LeaderStartSendAppendRequestTimeout(state *State)  {

	state.send_append_request_timeout_index ++
	index := state.send_append_request_timeout_index

	go func() {
		if state.User != "leader" { return }
		time.Sleep(time.Duration(1000) * time.Millisecond)

		if state.User != "leader" { return }
		if state.send_append_request_timeout_index != index { return }
		state.EventChan <- "send_append_request_timeout"
	}()
}

func LeaderSendAppendRequest(state *State)  {
	// log
	log.Info(state.Ip, ": ", "leader sent AppendRequest to ", state.Ips)

	go func() {
		for _, ip := range state.Ips {
			if state.User != "leader" { return }
			packet := Packet{Name: "AppendRequest", Id: state.Ip, DST_ID:ip, Term: state.CurrentTerm}
			state.SendPacketChan <- packet
		}
	}()
}
