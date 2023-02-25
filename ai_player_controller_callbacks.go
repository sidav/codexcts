package main

import "log"

func (ai *aiPlayerController) selectCoordsFromListCallback(message string, coords []*playerZoneCoords) *playerZoneCoords {
	coord := coords[rnd.Rand(len(coords))]
	log.Printf("  I chose coords: player %s, zone %d, index %d, which is %s \n", coord.player.name, coord.zone, coord.indexInZone, coord.getFormattedName())
	return coord
}

func (ai *aiPlayerController) showMessage(a, b string) {
	// stub
}
