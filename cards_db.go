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
	// NEUTRAL STARTER DECK
	&unitCard{
		cost:         1,
		name:         "Timely Messenger",
		subtype:      "Mercenary",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      1,
		baseHP:       1,
		specials:     []unitSpecial{{name: "Haste"}},
		startingDeck: true,
	},
	&unitCard{
		cost:         1,
		name:         "Tenderfoot",
		subtype:      "Virtuoso",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      1,
		baseHP:       2,
		startingDeck: true,
	},
	&unitCard{
		cost:         2,
		name:         "Older Brother",
		subtype:      "Drunkard",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      2,
		baseHP:       2,
		startingDeck: true,
	},
	&unitCard{
		cost:         2,
		name:         "Brick Thief",
		subtype:      "Mercenary",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      2,
		baseHP:       1,
		startingDeck: true,
	},
	&unitCard{
		cost:         2,
		name:         "Helpful Turtle",
		subtype:      "Cute Animal",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      1,
		baseHP:       2,
		specials:     []unitSpecial{{name: "Healing", value: 1}},
		startingDeck: true,
	},
	&unitCard{
		cost:         3,
		name:         "Granfalloon Flagbearer",
		subtype:      "Flagbearer",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      2,
		baseHP:       2,
		startingDeck: true,
	},
	&unitCard{
		cost:         3,
		name:         "Fruit Ninja",
		subtype:      "Ninja",
		element:      ELEMENT_NEUTRAL,
		techLevel:    0,
		baseAtk:      2,
		baseHP:       2,
		specials:     []unitSpecial{{name: "Frenzy", value: 1}},
		startingDeck: true,
	},
	&magicCard{
		cost:         1,
		name:         "Spark",
		subtype:      "Burn",
		isMinor:      true,
		description:  "Deal 1 damage to a patroller",
		element:      ELEMENT_NEUTRAL,
		startingDeck: true,
	},
	&magicCard{
		cost:         2,
		name:         "Bloom",
		subtype:      "Buff",
		isMinor:      true,
		description:  "Put +1/+1 rune on a friendly unit or hero that doesn't have +1/+1 rune",
		element:      ELEMENT_NEUTRAL,
		startingDeck: true,
	},
	&magicCard{
		cost:         2,
		name:         "Wither",
		subtype:      "Debuff",
		isMinor:      true,
		description:  "Put a -1/-1 rune on a unit or hero",
		element:      ELEMENT_NEUTRAL,
		startingDeck: true,
	},
	// BASHING DECK
	&magicCard{
		cost:         0,
		name:         "Wrecking Ball",
		subtype:      "",
		description:  "Deal 2 damage to a building",
		element:      ELEMENT_BASHING,
		startingDeck: true,
	},
	&magicCard{
		cost:        0,
		name:        "Wrecking Ball",
		subtype:     "Debuff",
		description: "Destroy a tech 0 or tech 1 unit.",
		element:     ELEMENT_BASHING,
	},
	&magicCard{
		cost:        1,
		name:        "Intimidate",
		subtype:     "Debuff",
		description: "Give a unit or hero -4 ATK this turn.",
		element:     ELEMENT_BASHING,
	},
	&magicCard{
		cost:        6,
		name:        "Final Smash",
		isUltimate:  true,
		subtype:     "Debuff",
		description: "Destroy a tech 0 unit, return a tech 1 unit to its owner's hand, and gain control of a tech 2 unit.",
		element:     ELEMENT_BASHING,
	},
	&unitCard{
		cost:      3,
		name:      "Iron Man",
		subtype:   "Mercenary",
		element:   ELEMENT_BASHING,
		techLevel: 1,
		baseAtk:   3,
		baseHP:    4,
		// specials:     []unitSpecial{{name: "Frenzy", value: 1}},
	},
	&unitCard{
		cost:      2,
		name:      "Revolver Ocelot",
		subtype:   "Leopard",
		element:   ELEMENT_BASHING,
		techLevel: 1,
		baseAtk:   3,
		baseHP:    3,
		specials:  []unitSpecial{{name: "Sparkshot"}},
	},
	&unitCard{
		cost:      4,
		name:      "Hired Stomper",
		subtype:   "Lizardman",
		element:   ELEMENT_BASHING,
		techLevel: 2,
		baseAtk:   4,
		baseHP:    3,
		// Arrives: deals 3 damage to a unit.
	},
	&unitCard{
		cost:      4,
		name:      "Regular-Sized Rhinoceros",
		subtype:   "Rhino",
		element:   ELEMENT_BASHING,
		techLevel: 2,
		baseAtk:   5,
		baseHP:    6,
	},
	&unitCard{
		cost:      3,
		name:      "Sneaky Pig",
		subtype:   "Pig",
		element:   ELEMENT_BASHING,
		techLevel: 2,
		baseAtk:   3,
		baseHP:    3,
		specials:  []unitSpecial{{name: "Haste"}},
		// Arrives: gets stealth this turn
	},
	&unitCard{
		cost:      4,
		name:      "Eggship",
		subtype:   "Contraption",
		element:   ELEMENT_BASHING,
		techLevel: 2,
		baseAtk:   4,
		baseHP:    3,
		specials:  []unitSpecial{{name: "Flying"}},
		// Arrives: gets stealth this turn
	},
	&unitCard{
		cost:      7,
		name:      "Trojan Duck",
		subtype:   "Contraption",
		element:   ELEMENT_BASHING,
		techLevel: 3,
		baseAtk:   8,
		baseHP:    9,
		specials:  []unitSpecial{{name: "Obliterate", value: 2}},
		// Arrives or attacks: deal 4 damage to a building
	},
}
