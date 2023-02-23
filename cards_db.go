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
		name:        "The Boot",
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
		cost:      5,
		name:      "Harvest Reaper",
		subtype:   "Contraption",
		element:   ELEMENT_BASHING,
		techLevel: 2,
		baseAtk:   6,
		baseHP:    5,
		specials:  []unitSpecial{{name: "Overpower"}},
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

	// FINESSE DECK
	&magicCard{
		cost:      2,
		name:      "Harmony",
		isOngoing: true,
		subtype:   "Buff",
		// channeling
		description: "Whenever you play a spell, summon a 0/1 neutral Dancer token (limit: 3). \n Sacrifice Harmony: " +
			"\"Stop the music.\" (Flip Dancers tokens)",
		element: ELEMENT_FINESSE,
	},
	&magicCard{
		cost:        2,
		name:        "Discord",
		subtype:     "Debuff",
		description: "Give all opponent's tech 0 and tech 1 units -2/-1 until end of turn.",
		element:     ELEMENT_FINESSE,
	},
	&magicCard{
		cost:      2,
		name:      "Two Step",
		isOngoing: true,
		subtype:   "Buff (target)",
		// channeling
		description: "Two of your units become dance partners if they aren't partnered already. While you control both, " +
			"they each get +2/+2. If you lose one, sacrifice Two Step.",
		element: ELEMENT_FINESSE,
	},
	&magicCard{
		cost:       1,
		name:       "Appel Stomp",
		isUltimate: true,
		subtype:    "Debuff",
		// channeling
		description: "Sideline a patroller, draw a card, then you may put Appel Stomp on top of your draw pile.",
		element:     ELEMENT_FINESSE,
	},
	&unitCard{
		cost:      2,
		name:      "Nimble Fencer",
		subtype:   "Virtuoso",
		element:   ELEMENT_FINESSE,
		techLevel: 1,
		baseAtk:   2,
		baseHP:    3,
		// specials:     []unitSpecial{{name: "Frenzy", value: 1}},
		// your Virtuosos have haste.
	},
	&unitCard{
		cost:      2,
		name:      "Star-Crossed Starlet",
		subtype:   "Virtuoso",
		element:   ELEMENT_FINESSE,
		techLevel: 1,
		baseAtk:   3,
		baseHP:    2,
		// specials:     []unitSpecial{{name: "Frenzy", value: 1}},
		// upkeep: this takes 1 damage. This gets +1 ATK for each damage on her.
	},
	&unitCard{
		cost:      5,
		name:      "Grounded Guide",
		subtype:   "Thespian",
		element:   ELEMENT_FINESSE,
		techLevel: 2,
		baseAtk:   4,
		baseHP:    4,
		// Your other units get +1 ATK. Your Virtuosos get +2/+1 instead.
	},
	&unitCard{
		cost:      3,
		name:      "Maestro",
		subtype:   "Thespian",
		element:   ELEMENT_FINESSE,
		techLevel: 2,
		baseAtk:   3,
		baseHP:    5,
		// Your Virtuosos cost 0 to play and gain "TAP: deal 2 damage to a building (target)"
		// specials:     []unitSpecial{{name: "Frenzy", value: 1}},
	},
	&unitCard{
		cost:      3,
		name:      "Backstabber",
		subtype:   "Rogue",
		element:   ELEMENT_FINESSE,
		techLevel: 2,
		baseAtk:   3,
		baseHP:    3,
		specials:  []unitSpecial{{name: "Invisible"}},
	},
	&unitCard{
		cost:      2,
		name:      "Cloud Sprite",
		subtype:   "Fairy",
		element:   ELEMENT_FINESSE,
		techLevel: 2,
		baseAtk:   3,
		baseHP:    2,
		specials:  []unitSpecial{{name: "Flying"}},
	},
	&unitCard{
		cost:      1,
		name:      "Leaping Lizard",
		subtype:   "Lizardman",
		element:   ELEMENT_FINESSE,
		techLevel: 2,
		baseAtk:   3,
		baseHP:    5,
		specials:  []unitSpecial{{name: "Anti-air"}},
	},
	&unitCard{
		cost:      6,
		name:      "Blademaster",
		subtype:   "Virtuoso",
		element:   ELEMENT_FINESSE,
		techLevel: 3,
		baseAtk:   7,
		baseHP:    5,
		// your units and heroes have swift strike.
	},
}
