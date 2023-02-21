package main

type card interface {
	getCost() int
	getName() string
	isInStartingDeck() bool
	getElement() element
}
