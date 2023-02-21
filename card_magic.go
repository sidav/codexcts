package main

type magicCard struct {
	cost         int
	name         string
	element      element
	startingDeck bool
}

func (mc *magicCard) getName() string {
	return mc.name
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
