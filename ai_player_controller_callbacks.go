package main

import "log"

func (ai *aiPlayerController) selectCoordsFromListCallback(message string, coords []*playerZoneCoords) *playerZoneCoords {
	index := rnd.SelectRandomIndexFromWeighted(len(coords), func(x int) int {
		switch coords[x].zone {
		case PLAYERZONE_PATROL, PLAYERZONE_OTHER:
			return 1
		default:
			return 5
		}
	})
	coord := coords[index]
	log.Printf("  I chose coords: player %s, zone %d, index %d, which is %s \n", coord.player.name, coord.zone, coord.indexInZone, coord.getFormattedName())
	return coord
}

func (ai *aiPlayerController) showMessage(a, b string) {
	// stub
}
