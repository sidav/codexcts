package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)

const (
	patrolZoneW = 70
	patrolZoneH = 10
	cardShortH  = 8
	cardFullW   = 30
	cardFullH   = 20
)

type tcellRenderer struct {
	w, h         int
	currUiLine   int
	g            *game
	activePlayer *player
	enemy        *player
	pc           *playerController
}

func (r *tcellRenderer) renderGame(g *game, renderForPlayerNum int, pc *playerController) {
	r.g = g
	r.activePlayer = g.players[renderForPlayerNum]
	r.enemy = g.players[(renderForPlayerNum+1)%2]
	r.w, r.h = cw.GetConsoleSize()
	cw.ClearScreen()

	r.renderHeader()
	r.renderEnemyField()
	r.renderPlayerField()

	r.renderPcmodeSpecific()

	cw.FlushScreen()
}

func (r *tcellRenderer) renderPcmodeSpecific() {
	switch pc.currentMode {
	case PCMODE_NONE:
	case PCMODE_CARD_FROM_HAND_SELECTED:
		r.renderSelectedCardFromHand()
	case PCMODE_UNIT_SELECTED:
		r.renderSelectedUnit()
	case PCMODE_MOVE_SELECTED_UNIT:
		cw.SetStyle(tcell.ColorBlack, tcell.ColorYellow)
		r.drawFilledInfoRect(fmt.Sprintf("Where to move %s?", pc.currentSelectedUnit.card.getName()), r.w/2, r.h-cardShortH)
	case PCMODE_SELECT_BUILDING:
		r.renderSelectBuildingMenu()
	case PCMODE_SELECT_HERO_TO_PLAY:
		r.renderCommandZone(r.activePlayer)
	default:
		panic("Check for pc mode specifics in renderer!")
	}
}

func (r *tcellRenderer) renderSelectBuildingMenu() {
	ww, wh := r.w/4, r.h/2
	wx, wy := r.w/2-ww/2, r.h/2-wh/2
	r.drawWindow("SELECT BUILDING", wx, wy, ww, wh, tcell.ColorBlue)
	r.currUiLine = wy + 2
	for _, b := range sTableBuildings {
		hotkey := "NO HOTKEY"
		if r.g.canPlayerBuild(r.activePlayer, b) {
			cw.ResetStyle()
		} else {
			cw.SetFg(tcell.ColorDarkGray)
		}
		if b.givesTech > 0 {
			hotkey = "T"
			r.drawLineAndIncrementY(fmt.Sprintf("%s - build %s", hotkey, b.name), wx+2)
		} else {
			switch b.name {
			case "Tower":
				hotkey = "O"
			case "Surplus":
				hotkey = "S"
			}
			r.drawLineAndIncrementY(fmt.Sprintf("%s - build %s (Add-on)", hotkey, b.name), wx+2)
		}
		r.drawLineAndIncrementY(fmt.Sprintf("   $%d, requires %d workers", b.cost, b.requiresWorkers), wx+1)
		r.currUiLine++
	}
}

func (r *tcellRenderer) renderHeader() {
	cw.SetStyle(tcell.ColorBlack, tcell.ColorYellow)
	cw.DrawRect(0, 0, r.w, 0)
	cw.SetStyle(tcell.ColorYellow, tcell.ColorBlack)
	cw.PutStringCenteredAt(fmt.Sprintf(" TURN %d: PLAYER %d - %s phase ",
		r.g.currentTurn, r.g.currentPlayerNumber, r.g.getCurrentPhaseName()), r.w/2, 0)
}

func (r *tcellRenderer) renderEnemyField() {
	r.currUiLine = 1
	r.renderOtherZone(r.enemy, 14, r.currUiLine)
	cw.SetStyle(tcell.ColorRed, tcell.ColorBlack)
	r.drawLineAndIncrementY(fmt.Sprintf("Base HP %d", r.enemy.baseHealth), 0)
	cw.ResetStyle()
	r.drawLineAndIncrementY(fmt.Sprintf("HAND: %4d", len(r.enemy.hand)), 0)
	r.drawLineAndIncrementY(fmt.Sprintf("DRAW: %4d", len(r.enemy.draw)), 0)
	r.drawLineAndIncrementY(fmt.Sprintf("DISCARD: %d", len(r.enemy.discard)), 0)
	r.drawLineAndIncrementY(fmt.Sprintf("WORKERS: %d", r.enemy.workers), 0)
	cw.SetFg(tcell.ColorYellow)
	r.drawLineAndIncrementY(fmt.Sprintf("$%d", r.enemy.gold), 0)

	for _, b := range r.enemy.techBuildings {
		if b != nil {
			r.drawLineAndIncrementY(b.static.name, 0)
		}
	}
	if r.enemy.addonBuilding != nil {
		r.drawLineAndIncrementY(r.enemy.addonBuilding.static.name, 0)
	}

	r.renderPatrolZone(r.enemy, 2)
}

func (r *tcellRenderer) renderPlayerField() {
	r.currUiLine = r.h/2 + 1
	r.renderOtherZone(r.activePlayer, 14, r.currUiLine)
	cw.SetStyle(tcell.ColorRed, tcell.ColorBlack)
	r.drawLineAndIncrementY(fmt.Sprintf("Base HP %d", r.activePlayer.baseHealth), 0)
	cw.ResetStyle()
	for _, tb := range r.activePlayer.techBuildings {
		if tb != nil {
			line := tb.static.name
			if tb.isUnderConstruction {
				line += " (under construction)"
			}
			r.drawLineAndIncrementY(line, 0)
		}
	}
	r.drawLineAndIncrementY("(B)build", 0)
	r.drawLineAndIncrementY(fmt.Sprintf("DRAW: %4d", len(r.activePlayer.draw)), 0)
	r.drawLineAndIncrementY(fmt.Sprintf("DISCARD: %d", len(r.activePlayer.discard)), 0)
	r.drawLineAndIncrementY(fmt.Sprintf("WORKERS: %d", r.activePlayer.workers), 0)
	cw.SetFg(tcell.ColorYellow)
	r.drawLineAndIncrementY(fmt.Sprintf("$%d", r.activePlayer.gold), 0)
	r.renderPatrolZone(r.activePlayer, r.h-cardShortH-patrolZoneH-2)
	r.renderHand()
}

func (r *tcellRenderer) renderHand() {
	cards := len(r.activePlayer.hand)
	cardW := r.w / cards
	if cardW > r.w/5 {
		cardW = r.w / 5
	}
	for i, c := range r.activePlayer.hand {
		r.renderCardInHand(c, i*(cardW), r.h-cardShortH, cardW, cardShortH)
	}
}

func (r *tcellRenderer) renderSelectedCardFromHand() {
	r.renderCardFull(pc.currentSelectedCardFromHand, r.w/2-cardFullW/2, r.h/2-cardFullH/2, cardFullW, cardFullH)
	cw.SetStyle(tcell.ColorBlack, tcell.ColorBlue)
	r.currUiLine = r.h/2 + cardFullH/2 + 1
	r.drawLineAndIncrementY(fmt.Sprintf(" %-30s", "W - spend as worker"), r.w/2-cardFullW/2)
	r.drawLineAndIncrementY(fmt.Sprintf(" %-30s", "P - play card"), r.w/2-cardFullW/2)
}

func (r *tcellRenderer) renderSelectedUnit() {
	r.renderCardFull(pc.currentSelectedUnit.card, r.w/2-cardFullW/2, r.h/2-cardFullH/2, cardFullW, cardFullH)
	r.currUiLine = r.h/2 + cardFullH/2 + 1
	cw.SetStyle(tcell.ColorBlack, tcell.ColorBlue)
	r.drawLineAndIncrementY(fmt.Sprintf(" %-30s", "M - move to other zone"), r.w/2-cardFullW/2)
	r.drawLineAndIncrementY(fmt.Sprintf(" %-30s", "A - attack"), r.w/2-cardFullW/2)
	r.drawLineAndIncrementY(fmt.Sprintf(" %-30s", "U - use ability"), r.w/2-cardFullW/2)
	r.drawLineAndIncrementY(fmt.Sprintf(" %-30s", "L - level up"), r.w/2-cardFullW/2)
}

func (r *tcellRenderer) renderOtherZone(p *player, x, y int) {
	const keys = "QWERT"
	for i, c := range p.otherZone {
		str := fmt.Sprintf("%s - %s",
			string(keys[i]),
			c.card.getName(),
		)
		cw.PutString(str, x, y+i)
	}
}

func (r *tcellRenderer) renderCommandZone(p *player) {
	ww, wh := r.w-r.w/5, r.h-r.h/10
	wx, wy := (r.w-ww)/2, (r.h-wh)/2
	r.drawWindow("COMMAND ZONE", wx, wy, ww, wh, tcell.ColorBlue)
	cardW, cardH := ww/3-4, wh-4
	cx, cy := wx+2, wy+2
	for i, h := range p.commandZone {
		if h != nil {
			r.renderCardFull(h, cx+i*(cardW+2), cy, cardW, cardH)
		}
	}
}

func (r *tcellRenderer) renderPatrolZone(p *player, y int) {
	x := r.w - patrolZoneW - 1
	cw.ResetStyle()
	cardW := patrolZoneW / 5
	for i := 0; i < 5; i++ {
		currX := x + i*cardW
		cw.SetStyle(tcell.ColorBlack, tcell.ColorGray)
		cw.DrawRect(currX, y, cardW, patrolZoneH)
		cw.SetStyle(tcell.ColorGray, tcell.ColorBlack)
		descrString := ""
		hotkey := ""
		switch i {
		case 0:
			cw.PutStringCenteredAt(" Squad leader", currX+cardW/2, y+1)
			hotkey = "Y"
			descrString = "+SHLD/Taunt"
		case 1:
			cw.PutStringCenteredAt("Elite", currX+cardW/2, y+1)
			hotkey = "U"
			descrString = "+1 ATK"
		case 2:
			cw.PutStringCenteredAt("Scavenger", x+i*cardW+cardW/2, y+1)
			hotkey = "I"
			descrString = "Dies: +1$"
		case 3:
			cw.PutStringCenteredAt("Technician", x+i*cardW+cardW/2, y+1)
			hotkey = "O"
			descrString = "Dies: +1 Card"
		case 4:
			cw.PutStringCenteredAt("Lookout", x+i*cardW+cardW/2, y+1)
			hotkey = "P"
			descrString = "Resist 1"
		}
		cw.PutStringCenteredAt(descrString, currX+cardW/2, y+patrolZoneH)
		unitHere := p.patrolZone[i]
		if unitHere != nil {
			r.drawUnit(unitHere, x+i*cardW+1, y+1, cardW-1, patrolZoneH-2)
		}
		if unitHere == nil && p == r.activePlayer {
			cw.PutStringCenteredAt(hotkey, currX+cardW/2, y+patrolZoneH/2)
		}
	}
}

func (r *tcellRenderer) drawLineAndIncrementY(line string, x int) {
	cw.PutString(line, x, r.currUiLine)
	r.currUiLine++
}

func (r *tcellRenderer) drawWindow(header string, x, y, w, h int, borderColor tcell.Color) {
	cw.SetStyle(tcell.ColorBlack, borderColor)
	cw.DrawRect(x, y, w, h)
	cw.SetStyle(borderColor, tcell.ColorBlack)
	cw.PutStringCenteredAt(" "+header+" ", x+w/2, y)
	cw.DrawFilledRect(' ', x+1, y+1, w-2, h-2)
	cw.ResetStyle()
}

func (r *tcellRenderer) drawFilledInfoRect(text string, centerX, centerY int) {
	w := len(text) + 2
	h := 3
	cw.DrawFilledRect(' ', centerX-w/2, centerY-1, w, h-1)
	cw.PutStringCenteredAt(text, centerX, centerY)
}
