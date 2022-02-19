package main

type state int

const (
	DEFAULT state = iota
	AUTH
	CLASSIFY
	ADD_ACCOUNT
	BANKS
)
