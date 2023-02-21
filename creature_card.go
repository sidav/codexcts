package main

type creatureCard struct {
	cost      int
	name      string
	element   uint8
	techLevel int
	baseAtk   int
	baseDef   int
}

func (cc *creatureCard) getCost() int {
	return cc.cost
}
