package main

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

func (cc *creatureCard) getCost() int {
	return cc.cost
}

func (cc *creatureCard) isInStartingDeck() bool {
	return cc.startingDeck
}

func (cc *creatureCard) getElement() element {
	return cc.element
}
