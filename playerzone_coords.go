package main

import "fmt"

const (
	PLAYERZONE_OTHER = iota
	PLAYERZONE_PATROL
	PLAYERZONE_TECH_BUILDINGS
	PLAYERZONE_ADDON_BUILDING
	PLAYERZONE_MAIN_BASE
)

type playerZoneCoords struct {
	player      *player
	zone        int
	indexInZone int
}

func (pzc *playerZoneCoords) getFormattedName() string {
	switch pzc.zone {
	case PLAYERZONE_MAIN_BASE:
		return fmt.Sprintf("%s's Main Base", pzc.player.name)
	case PLAYERZONE_TECH_BUILDINGS:
		return fmt.Sprintf("%s's %s", pzc.player.name, pzc.player.techBuildings[pzc.indexInZone].static.name)
	case PLAYERZONE_ADDON_BUILDING:
		return fmt.Sprintf("%s's %s", pzc.player.name, pzc.player.addonBuilding.static.name)
	case PLAYERZONE_OTHER:
		return fmt.Sprintf("%s's %s in other zone", pzc.player.name, pzc.player.otherZone[pzc.indexInZone].getName())
	case PLAYERZONE_PATROL:
		return fmt.Sprintf("%s's %s in patrol zone (slot %d)", pzc.player.name,
			pzc.player.patrolZone[pzc.indexInZone].getName(), pzc.indexInZone)
	default:
		panic("Coords formatting error.")
	}
}
