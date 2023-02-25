package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)

func (r *tcellRenderer) renderAttackSelection() {
	length := len(r.pc.callbackCoordsList)
	ww, wh := r.w/3, length+2
	wx, wy := r.w/2-ww/2, r.h/2-wh/2
	r.drawWindow(r.pc.callbackMessage, wx, wy, ww, wh, tcell.ColorRed)
	r.currUiLine = wy + 1
	for index, coord := range r.pc.callbackCoordsList {
		coordPlayer := coord.player
		switch coord.zone {
		case PLAYERZONE_MAIN_BASE:
			r.drawLineAndIncrementY(fmt.Sprintf("%s - %s (%d HP)", string(coordsSelectableWithCallback[index]),
				"Main base", coordPlayer.baseHealth), wx+1)
		case PLAYERZONE_TECH_BUILDINGS:
			r.drawLineAndIncrementY(fmt.Sprintf("%s - %s (%d HP)", string(coordsSelectableWithCallback[index]),
				coordPlayer.techBuildings[coord.indexInZone].static.name,
				coordPlayer.techBuildings[coord.indexInZone].currentHitpoints,
			), wx+1)
		case PLAYERZONE_ADDON_BUILDING:
			r.drawLineAndIncrementY(fmt.Sprintf("%s - %s (%d HP)", string(coordsSelectableWithCallback[index]),
				coordPlayer.addonBuilding.static.name,
				coordPlayer.addonBuilding.currentHitpoints,
			), wx+1)
		case PLAYERZONE_OTHER:
			unt := coordPlayer.otherZone[coord.indexInZone]
			_, untHp := unt.getAtkHp()
			r.drawLineAndIncrementY(fmt.Sprintf("%s - %s (%d HP)", string(coordsSelectableWithCallback[index]),
				unt.card.getName(), untHp), wx+1)
		case PLAYERZONE_PATROL:
			unt := coordPlayer.patrolZone[coord.indexInZone]
			untAtk, untHp := unt.getAtkHp()
			r.drawLineAndIncrementY(fmt.Sprintf("%s - %d/%d %s", string(coordsSelectableWithCallback[index]),
				untAtk, untHp, unt.card.getName()), wx+1)
		}
	}
}
