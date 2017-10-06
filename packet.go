package main

type Packet struct {
	Name string
	Id string
	Term int
	DST_ID string

	// election
	VoteGranted bool

	// log replication
}
