package main

type Packet struct {
	Name string
	Id string
	Term int

	// election
	VoteGranted bool

	// log replication
}
