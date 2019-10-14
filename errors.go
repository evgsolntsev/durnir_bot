package main

type Error struct {
	text string
}

func (e Error) Error() string {
	return e.text
}

func (e Error) String() string {
	return e.text
}
