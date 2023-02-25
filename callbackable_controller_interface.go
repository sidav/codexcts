package main

type callbackableController interface {
	act(g *game)
	phaseEnded() bool
	selectCoordsFromListCallback(string, []*playerZoneCoords) *playerZoneCoords
}
