package main

import "log"

func (ai *aiPlayerController) selectCoordsFromListCallback(message string, coords []*playerZoneCoords) *playerZoneCoords {
	coord := coords[rnd.Rand(len(coords))]
	log.Printf("  I chose coords: player %s, zone %d, index %d", coord.player.name, coord.zone, coord.indexInZone)
	return coord
}
