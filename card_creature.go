package main

import "fmt"

type creatureCard struct {
	cost         int
	name         string
	element      element
	techLevel    int
	baseAtk      int
	baseDef      int
	startingDeck bool
}

func (cc *creatureCard) getName() string {
	return cc.name
}

func (cc *creatureCard) getFormattedName() string {
	return fmt.Sprintf("(%d) %-25s %d/%d", cc.cost, cc.name, cc.baseAtk, cc.baseDef)
}

func (cc *creatureCard) getCost() int {
	return cc.cost
}

func (cc *creatureCard) isInStartingDeck() bool {
	return cc.startingDeck
}

func (cc *creatureCard) getElement() element {
	return cc.element
}
