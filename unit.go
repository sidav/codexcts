package main

// unit is "card on the battlefield", as opposed to "card in deck"
type unit struct {
	card   card
	tapped bool
	wounds int
	level  int // needed only for heroes
}

func (u *unit) getAtkHp() (int, int) {
	switch u.card.(type) {
	case *unitCard:
		return u.card.(*unitCard).baseAtk, u.card.(*unitCard).baseHP
	}
	return 0, 0
}
