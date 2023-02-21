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
	&creatureCard{
		cost:         1,
		name:         "Timely Messenger",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      1,
		baseDef:      1,
		startingDeck: true,
	},
	&creatureCard{
		cost:         1,
		name:         "Tenderfoot",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      1,
		baseDef:      2,
		startingDeck: true,
	},
	&creatureCard{
		cost:         2,
		name:         "Older Brother",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      2,
		baseDef:      2,
		startingDeck: true,
	},
	&creatureCard{
		cost:         2,
		name:         "Brick Thief",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      2,
		baseDef:      1,
		startingDeck: true,
	},
	&creatureCard{
		cost:         2,
		name:         "Helpful Turtle",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      1,
		baseDef:      2,
		startingDeck: true,
	},
	&creatureCard{
		cost:         3,
		name:         "Granfalloon Flagbearer",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      2,
		baseDef:      2,
		startingDeck: true,
	},
	&creatureCard{
		cost:         3,
		name:         "Fruit Ninja",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      2,
		baseDef:      2,
		startingDeck: true,
	},
	&magicCard{
		cost:         1,
		name:         "Spark",
		element:      ELEMENT_NEUTRAL,
		startingDeck: true,
	},
	&magicCard{
		cost:         2,
		name:         "Bloom",
		element:      ELEMENT_NEUTRAL,
		startingDeck: true,
	},
	&magicCard{
		cost:         2,
		name:         "Wither",
		element:      ELEMENT_NEUTRAL,
		startingDeck: true,
	},
}
