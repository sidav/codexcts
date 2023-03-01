package main

import (
	"github.com/gdamore/tcell/v2"
	"strings"
)

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
		if index != -1 && index < len(coords) {
			pc.currentMode = previousMode
			return coords[index]
		}
	}
}

func (pc *playerController) showMessage(title, msg string) {
	io.showMessageWindow(title, msg, tcell.ColorRed)
	for readKey() != "ENTER" {
	}
}
