package state

type State int

const (
	HOME State = iota
	AUTH
	CLASSIFY
	SYNC
)
