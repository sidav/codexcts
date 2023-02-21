package main

import "fmt"

type heroCard struct {
	cost         int
	name         string
	element      element
	levelsAttDef [][3]int // [level baseAtt baseDef]
}

func (ch *heroCard) getName() string {
	return ch.name
}

func (ch *heroCard) getFormattedName() string {
	return fmt.Sprintf("(%d) %-25s", ch.cost, ch.name)
}

func (ch *heroCard) getCost() int {
	return ch.cost
}
