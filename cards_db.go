package main

func getStartingCardsForElement(elem element) []card {
	cards := make([]card, 0)
	for _, c := range cardsDb {
		if c.isInStartingDeck() && c.getElement() == elem {
			cards = append(cards, c)
		}
	}
	return cards
}

func getCardByName(name string) card {
	for _, c := range cardsDb {
		if c.getName() == name {
			return c
		}
	}
	panic("No such card in DB: " + name)
}

var cardsDb = []card{
	&unitCard{
		cost:         1,
		name:         "Timely Messenger",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      1,
		baseDef:      1,
		specials:     []unitSpecial{{name: "Haste"}},
		startingDeck: true,
	},
	&unitCard{
		cost:         1,
		name:         "Tenderfoot",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      1,
		baseDef:      2,
		startingDeck: true,
	},
	&unitCard{
		cost:         2,
		name:         "Older Brother",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      2,
		baseDef:      2,
		startingDeck: true,
	},
	&unitCard{
		cost:         2,
		name:         "Brick Thief",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      2,
		baseDef:      1,
		startingDeck: true,
	},
	&unitCard{
		cost:         2,
		name:         "Helpful Turtle",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      1,
		baseDef:      2,
		specials:     []unitSpecial{{name: "Healing", value: 1}},
		startingDeck: true,
	},
	&unitCard{
		cost:         3,
		name:         "Granfalloon Flagbearer",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      2,
		baseDef:      2,
		startingDeck: true,
	},
	&unitCard{
		cost:         3,
		name:         "Fruit Ninja",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      2,
		baseDef:      2,
		specials:     []unitSpecial{{name: "Frenzy", value: 1}},
		startingDeck: true,
	},
	&magicCard{
		cost:         1,
		name:         "Spark",
		description:  "Deal 1 damage to a patroller",
		element:      ELEMENT_NEUTRAL,
		startingDeck: true,
	},
	&magicCard{
		cost:         2,
		name:         "Bloom",
		description:  "Put +1/+1 rune on a friendly unit or hero that doesn't have +1/+1 rune",
		element:      ELEMENT_NEUTRAL,
		startingDeck: true,
	},
	&magicCard{
		cost:         2,
		name:         "Wither",
		description:  "Put a -1/-1 rune on a unit or hero",
		element:      ELEMENT_NEUTRAL,
		startingDeck: true,
	},
}
