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
	},
}
