package main

// creature is "card on the battlefield", as opposed to "card in deck"
type creature struct {
	card   card
	wounds int
	level  int // needed only for heroes
}
