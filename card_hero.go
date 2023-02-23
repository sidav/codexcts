package main

import "fmt"

type heroCard struct {
	cost                 int
	name                 string
	element              element
	levelsAttDef         [][3]int // [level baseAtt baseDef]
	levelsAbilitiesTexts []string
}

func (ch *heroCard) getName() string {
	return ch.name
}

func (ch *heroCard) getElement() element {
	return ch.element
}

func (ch *heroCard) isInStartingDeck() bool {
	return false
}

func (ch *heroCard) getFormattedName() string {
	return fmt.Sprintf("(%d) %-25s", ch.cost, ch.name)
}

func (ch *heroCard) getCost() int {
	return ch.cost
}

func (ch *heroCard) getMaxLevel() int {
	return ch.levelsAttDef[len(ch.levelsAttDef)-1][0]
}
