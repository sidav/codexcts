package main

import "fmt"

// unit is "card on the battlefield", as opposed to "card in deck"
type unit struct {
	card             card
	tapped           bool // "exhausted" in rules
	attackedThisTurn bool // for units with Readiness, "tapped" is not enough
	wounds           int
	level            int // needed only for heroes
}

func (u *unit) getName() string {
	return u.card.getName()
}

func (u *unit) getNameWithStats() string {
	a, h := u.getAtkHp()
	return fmt.Sprintf("%d/%d %s", a, h, u.card.getName())
}

func (u *unit) isHero() bool {
	_, ok := u.card.(*heroCard)
	return ok
}

func (u *unit) getAtkHp() (int, int) {
	switch u.card.(type) {
	case *unitCard:
		return u.card.(*unitCard).baseAtk, u.card.(*unitCard).baseHP - u.wounds
	case *heroCard:
		hc := u.card.(*heroCard)
		for i := len(hc.levelsAttDef) - 1; i >= 0; i-- {
			if u.level >= hc.levelsAttDef[i][0] {
				return hc.levelsAttDef[i][1], hc.levelsAttDef[i][2] - u.wounds
			}
		}
	}
	return -99, -99
}

func (u *unit) hasPassiveAbility(code UPACode) bool {
	switch u.card.(type) {
	case *unitCard:
		uc := u.card.(*unitCard)
		for i := range uc.passiveAbilities {
			if uc.passiveAbilities[i].code == code {
				return true
			}
		}
	case *heroCard:
		hc := u.card.(*heroCard)
		for i := range hc.levelsPassiveAbilities {
			if u.level >= hc.levelsPassiveAbilities[i].availableFromLevel && hc.levelsPassiveAbilities[i].code == code {
				return true
			}
		}
	}
	return false
}

func (u *unit) getListOfPassiveAbilities() []unitPassiveAbility {
	switch u.card.(type) {
	case *unitCard:
		return u.card.(*unitCard).passiveAbilities
	case *heroCard:
		abs := make([]unitPassiveAbility, 0)
		hc := u.card.(*heroCard)
		for i := range hc.levelsPassiveAbilities {
			if u.level >= hc.levelsPassiveAbilities[i].availableFromLevel {
				abs = append(abs, hc.levelsPassiveAbilities[i])
			}
		}
		return abs
	}
	panic("getListOfPassiveAbilities()")
}

func (u *unit) getPassiveAbilityValue(code UPACode) int {
	sum := 0
	switch u.card.(type) {
	case *unitCard:
		uc := u.card.(*unitCard)
		for i := range uc.passiveAbilities {
			if uc.passiveAbilities[i].code == code {
				sum += uc.passiveAbilities[i].value
			}
		}
	case *heroCard:
		hc := u.card.(*heroCard)
		for i := range hc.levelsPassiveAbilities {
			if u.level >= hc.levelsPassiveAbilities[i].availableFromLevel && hc.levelsPassiveAbilities[i].code == code {
				sum += hc.levelsPassiveAbilities[i].value
			}
		}
	}
	return sum
}
