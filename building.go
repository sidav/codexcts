package main

type building struct {
	static              *buildingStatic
	currentHitpoints    int
	isUnderConstruction bool
}

type buildingStatic struct {
	name            string
	cost            int
	maxHitpoints    int
	requiresWorkers int
	givesTech       int
	isAddon         bool
}

func getBuildingStaticByName(name string) *buildingStatic {
	for i := range sTableBuildings {
		if sTableBuildings[i].name == name {
			return sTableBuildings[i]
		}
	}
	panic("No building with that name!")
}

func getTechBuildingByTechLevel(level int) *buildingStatic {
	for i := range sTableBuildings {
		if sTableBuildings[i].givesTech == level {
			return sTableBuildings[i]
		}
	}
	panic("No building with that name!")
}

var sTableBuildings = []*buildingStatic{
	{
		name:            "Tech 1",
		cost:            1,
		requiresWorkers: 6,
		maxHitpoints:    5,
		givesTech:       1,
		isAddon:         false,
	},
	{
		name:            "Tech 2",
		cost:            4,
		requiresWorkers: 8,
		maxHitpoints:    5,
		givesTech:       2,
		isAddon:         false,
	},
	{
		name:            "Tech 3",
		cost:            5,
		requiresWorkers: 10,
		maxHitpoints:    5,
		givesTech:       3,
		isAddon:         false,
	},
	{
		name:            "Tower",
		cost:            3,
		maxHitpoints:    4,
		requiresWorkers: 0,
		givesTech:       0,
		isAddon:         true,
	},
	{
		name:            "Surplus",
		cost:            5,
		maxHitpoints:    4,
		requiresWorkers: 0,
		givesTech:       0,
		isAddon:         true,
	},
	//{
	//	name:            "Tech Lab",
	//	maxHitpoints:    5,
	//	requiresWorkers: 0,
	//	givesTech:       0,
	//	isAddon:         true,
	//},
	//{
	//	name:            "Heroes Hall",
	//	maxHitpoints:    5,
	//	requiresWorkers: 0,
	//	givesTech:       0,
	//	isAddon:         true,
	//},
}
