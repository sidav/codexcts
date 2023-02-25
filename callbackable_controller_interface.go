package main

type callbackableController interface {
	selectCoordsFromListCallback(string, []*playerZoneCoords) *playerZoneCoords
}
