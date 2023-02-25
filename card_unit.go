package main

import "fmt"

type unitCard struct {
	cost             int
	name             string
	subtype          string
	element          element
	techLevel        int
	baseAtk          int
	baseHP           int
	passiveAbilities []unitPassiveAbility
	startingDeck     bool
}

func (uc *unitCard) getName() string {
	return uc.name
}

func (uc *unitCard) getSubtype() string {
	return uc.subtype
}

func (uc *unitCard) getFormattedName() string {
	return fmt.Sprintf("%-15s %d/%d", uc.name, uc.baseAtk, uc.baseHP)
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
