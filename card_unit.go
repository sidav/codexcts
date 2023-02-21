package main

import "fmt"

type unitCard struct {
	cost         int
	name         string
	element      element
	techLevel    int
	baseAtk      int
	baseDef      int
	specials     []unitSpecial
	startingDeck bool
}

func (uc *unitCard) getName() string {
	return uc.name
}

func (uc *unitCard) getFormattedName() string {
	return fmt.Sprintf("(%d) %-25s %d/%d", uc.cost, uc.name, uc.baseAtk, uc.baseDef)
}

func (uc *unitCard) getCost() int {
	return uc.cost
}

func (uc *unitCard) isInStartingDeck() bool {
	return uc.startingDeck
}

func (uc *unitCard) getElement() element {
	return uc.element
}

func (uc *unitCard) hasSpecial(s string) bool {
	for i := range uc.specials {
		if uc.specials[i].name == s {
			return true
		}
	}
	return false
}
