package main

var heroCardsDb = []*heroCard{
	{
		cost:    2,
		name:    "Trog Bashar",
		element: ELEMENT_BASHING,
		levelsAttDef: [][3]int{
			{1, 2, 3},
			{5, 3, 4},
			{8, 4, 5},
		},
		levelsPassiveAbilities: []unitPassiveAbility{
			{
				code:               UPA_READINESS,
				availableFromLevel: 8,
			},
		},
		levelsAbilitiesTexts: []string{
			"",
			"Attacks: deal 1 damage to that opponent's base (target).",
			"Readiness",
		},
	},
	{
		cost:    2,
		name:    "River Montoya",
		element: ELEMENT_FINESSE,
		levelsAttDef: [][3]int{
			{1, 2, 3},
			{3, 2, 4},
			{5, 3, 4},
		},
		levelsAbilitiesTexts: []string{
			"",
			"TAP: Sideline a tech 0 or tech 1 patroller.",
			"Your tech 0 units cost 1 less to play.",
		},
	},
}
