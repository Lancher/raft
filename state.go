package main


type State struct {
	User string
	CurrentTerm int

	Ip string
	Ips [] string

	Vote int
}

func (state *State) UpdateCurrentTerm(term int)  {
	state.CurrentTerm = term
	state.Vote = 1
}
