package main

func (ai *aiPlayerController) selectCoordsFromListCallback(message string, coords []*playerZoneCoords) *playerZoneCoords {
	return coords[rnd.Rand(len(coords))]
}
