package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)

func (r *tcellRenderer) showMessageWindow(title string, msg string, color tcell.Color) {
	ww, wh := r.w/2, 5*r.h/6
	wx, wy := r.w/2-ww/2, r.h/2-wh/2
	r.drawWindow(title, wx, wy, ww, wh, color)
	cw.ResetStyle()
	cw.PutTextInRect(msg, wx+1, wy+1, ww-2)
	cw.FlushScreen()
}

func (r *tcellRenderer) renderAttackSelection() {
	length := len(r.pc.callbackCoordsList)
	ww, wh := r.w/3, length+5
	wx, wy := r.w/2-ww/2, r.h/2-wh/2
	r.drawWindow(r.pc.callbackMessage, wx, wy, ww, wh, tcell.ColorRed)
	r.currUiLine = wy + 1
	zoneString := ""
	prevZoneString := ""
	for index, coord := range r.pc.callbackCoordsList {
		coordPlayer := coord.player
		var unt *unit
		selectionString := ""

		switch coord.zone {
		case PLAYERZONE_MAIN_BASE:
			zoneString = "BASE"
			selectionString = fmt.Sprintf("%s - %s (%d HP)", string(coordsSelectableWithCallback[index]),
				"Main base", coordPlayer.baseHealth)
		case PLAYERZONE_TECH_BUILDINGS:
			zoneString = "BASE"
			selectionString = fmt.Sprintf("%s - %s (%d HP)", string(coordsSelectableWithCallback[index]),
				coordPlayer.techBuildings[coord.indexInZone].static.name,
				coordPlayer.techBuildings[coord.indexInZone].currentHitpoints,
			)
		case PLAYERZONE_ADDON_BUILDING:
			zoneString = "BASE"
			selectionString = fmt.Sprintf("%s - %s (%d HP)", string(coordsSelectableWithCallback[index]),
				coordPlayer.addonBuilding.static.name,
				coordPlayer.addonBuilding.currentHitpoints,
			)
		case PLAYERZONE_OTHER:
			zoneString = "OTHER ZONE"
			unt = coordPlayer.otherZone[coord.indexInZone]
		case PLAYERZONE_PATROL:
			zoneString = "PATROL ZONE"
			unt = coordPlayer.patrolZone[coord.indexInZone]
		}
		if unt != nil {
			untAtk, untHp := unt.getAtkHp()
			selectionString = fmt.Sprintf("%s - %d/%d %s", string(coordsSelectableWithCallback[index]),
				untAtk, untHp, unt.getName())
		}
		// separator
		if zoneString != prevZoneString {
			cw.SetFg(tcell.ColorDarkMagenta)
			r.drawLineAndIncrementY(zoneString, wx+3)
		}
		cw.ResetStyle()
		r.drawLineAndIncrementY(selectionString, wx+1)
		prevZoneString = zoneString
	}
}
