package main

import (
	"time"
	log "github.com/sirupsen/logrus"
)

func FollowerHandleReceivePacket(state *State, rec Packet)  {
	// AppendRequest
	if rec.Name == "AppendRequest" {
		if state.CurrentTerm < rec.Term {
			// log
			log.Info(state.Ip, ": ", "follower got AppendRequest with greater term (normal)")

			// user, term, vote
			state.User = "follower"
			state.CurrentTerm = rec.Term
			state.Vote = 1

			// send packet
			packet := Packet{Name: "AppendResponse", Id: state.Ip, DST_ID: rec.Id, Term: state.CurrentTerm}
			state.SendPacketChan <- packet

			// start timeout
			FollowerStartWaitAppendRequestTimeout(state)
		} else if state.CurrentTerm == rec.Term {
			// log
			log.Info(state.Ip, ": ", "follower got AppendRequest with the same term (normal)")

			// send packet
			packet := Packet{Name: "AppendResponse", Id: state.Ip, DST_ID: rec.Id, Term: state.CurrentTerm}
			state.SendPacketChan <- packet

			// start timeout
			FollowerStartWaitAppendRequestTimeout(state)
		} else if state.CurrentTerm > rec.Term {
			// log
			log.Info(state.Ip, ": ", "follower got AppendRequest with the smaller term (ignore the request)")

			// send packet
			packet := Packet{Name: "AppendResponse", Id: state.Ip, DST_ID: rec.Id, Term: state.CurrentTerm}
			state.SendPacketChan <- packet
		}

	// AppendResponse
	} else if rec.Name == "AppendResponse" {
		if state.CurrentTerm < rec.Term {
			// log
			log.Warning(state.Ip, ": ", "follower got AppendResponse with greater term (should not happened)")
		} else if state.CurrentTerm == rec.Term {
			// log
			log.Warning(state.Ip, ": ", "follower got AppendResponse with the same term (should not happened)")
		} else if state.CurrentTerm > rec.Term {
			// log
			log.Info(state.Ip, ": ", "follower got AppendResponse with less term (ignore late response)")
		}

		// VoteRequest
	}else if rec.Name == "VoteRequest" {
		if state.CurrentTerm < rec.Term {
			// log
			log.Info(state.Ip, ": ", "follower got VoteRequest with greater term and then become candidate (normal)")

			// user, term, vote
			state.User = "follower"
			state.CurrentTerm = rec.Term
			state.Vote = 0

			// send packet
			packet := Packet{Name: "VoteResponse", Id: state.Ip, DST_ID: rec.Id, Term: state.CurrentTerm, VoteGranted:true}
			state.SendPacketChan <- packet
		} else if state.CurrentTerm == rec.Term {
			// log
			log.Info(state.Ip, ": ", "follower got VoteRequest with the same term (normal)")

			// send packet
			if state.Vote == 1 {
				state.Vote = 0
				packet := Packet{Name: "VoteResponse", Id: state.Ip, DST_ID: rec.Id, Term: state.CurrentTerm, VoteGranted:true}
				state.SendPacketChan <- packet
			} else {
				packet := Packet{Name: "VoteResponse", Id: state.Ip, DST_ID: rec.Id, Term: state.CurrentTerm, VoteGranted:false}
				state.SendPacketChan <- packet
			}
		} else if state.CurrentTerm > rec.Term {
			// log
			log.Info(state.Ip, ": ", "follower got VoteRequest with less term (normal)")

			// send packet
			packet := Packet{Name: "VoteResponse", Id: state.Ip, DST_ID: rec.Id, Term: state.CurrentTerm, VoteGranted:false}
			state.SendPacketChan <- packet
		}

	// VoteResponse
	}else if rec.Name == "VoteResponse" {
		if state.CurrentTerm < rec.Term {
			// log
			log.Warning(state.Ip, ": ", "follower got VoteResponse with greater term (should not happened)")
		} else if state.CurrentTerm == rec.Term {
			// log
			log.Warning(state.Ip, ": ", "follower got VoteResponse with the same term (should not happened)")
		} else if state.CurrentTerm > rec.Term {
			// log
			log.Info(state.Ip, ": ", "follower got VoteResponse with smaller term (ignore late response)")
		}
	}
}

func FollowerStartWaitAppendRequestTimeout(state *State)  {

	state.wait_append_request_timeout_index ++
	index := state.wait_append_request_timeout_index

	go func() {

		if state.User != "follower" { return }
		time.Sleep(1500 * time.Millisecond)

		if state.User != "follower" { return }
		if state.wait_append_request_timeout_index != index { return }
		state.EventChan <- "wait_append_request_timeout"
	}()
}

func FollowerBecomeCandidate(state *State)  {
	// user, term, vote
	state.User = "candidate"
	state.CurrentTerm ++
	state.Vote = 1

	// event
	state.EventChan <- "become_candidate"
}
