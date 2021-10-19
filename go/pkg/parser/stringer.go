package parser

type stringer string

func (s stringer) String() string { return string(s) + "\n" }
