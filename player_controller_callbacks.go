package main

import "strings"

const coordsSelectableWithCallback = "1234567890qwertyuiop"

func (pc *playerController) selectCoordsFromListCallback(message string, coords []*playerZoneCoords) *playerZoneCoords {
	previousMode := pc.currentMode
	pc.currentMode = PCMODE_CALLBACK_TARGET_SELECTION
	pc.callbackMessage = message
	pc.callbackCoordsList = coords
	io.renderGame(pc.g, pc.g.currentPlayerNumber, pc)
	for {
		key := readKey()
		index := strings.Index(coordsSelectableWithCallback, key)
		if index != -1 {
			pc.currentMode = previousMode
			return coords[index]
		}
	}
}
