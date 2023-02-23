package main

import "fmt"

type magicCard struct {
	cost         int
	name         string
	subtype      string // buff, debuff, etc
	isMinor      bool
	isUltimate   bool
	isOngoing    bool
	description  string
	element      element
	startingDeck bool
}

func (mc *magicCard) getName() string {
	return mc.name
}

func (mc *magicCard) getSubtype() string {
	return mc.subtype
}

func (mc *magicCard) getFormattedName() string {
	return fmt.Sprintf("(%d) %-25s", mc.cost, mc.name)
}

func (mc *magicCard) getCost() int {
	return mc.cost
}

func (mc *magicCard) isInStartingDeck() bool {
	return mc.startingDeck
}

func (mc *magicCard) getElement() element {
	return mc.element
}
