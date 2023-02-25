package main

func (pc *playerController) selectCoordsFromListCallback(message string, coords []*playerZoneCoords) *playerZoneCoords {

	return coords[rnd.Rand(len(coords))]
}
