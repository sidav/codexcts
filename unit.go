package main

// unit is "card on the battlefield", as opposed to "card in deck"
type unit struct {
	card   card
	tapped bool
	wounds int
	level  int // needed only for heroes
}