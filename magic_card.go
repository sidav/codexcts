package main

type magicCard struct {
	cost int
}

func (mc *magicCard) getCost() int {
	return mc.cost
}
