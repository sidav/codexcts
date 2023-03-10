package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"strconv"
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

func (r *tcellRenderer) updateBounds() {
	r.w, r.h = cw.GetConsoleSize()
	cw.ClearScreen()
}

func (r *tcellRenderer) renderGame(g *game, renderForPlayerNum int, pc *playerController) {
	r.updateBounds()
	r.pc = pc
	r.g = g
	r.activePlayer = g.players[renderForPlayerNum]
	r.enemy = g.players[(renderForPlayerNum+1)%2]

	r.renderHeader()
	r.renderEnemyField()
	r.renderPlayerField()

	r.renderPcmodeSpecific()

	if r.g.currentPhase == PHASE_CODEX {
		r.renderCodexSelection(r.activePlayer)
	}

	cw.FlushScreen()
}

func (r *tcellRenderer) renderPcmodeSpecific() {
	switch r.pc.currentMode {
	case PCMODE_NONE:
	case PCMODE_CARD_FROM_HAND_SELECTED:
		r.renderSelectedCardFromHand()
	case PCMODE_UNIT_SELECTED:
		r.renderSelectedUnit()
	case PCMODE_MOVE_SELECTED_UNIT:
		cw.SetStyle(tcell.ColorBlack, tcell.ColorYellow)
		r.drawFilledInfoRect(fmt.Sprintf("Where to move %s?", r.pc.currentSelectedUnit.getName()), r.w/2, r.h-cardShortH)
	case PCMODE_SELECT_BUILDING:
		r.renderSelectBuildingMenu()
	case PCMODE_SELECT_HERO_TO_PLAY:
		r.renderCommandZone(r.activePlayer)
	case PCMODE_CALLBACK_TARGET_SELECTION:
		r.renderAttackSelection()
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

	r.renderPatrolZone(r.enemy, r.h-cardShortH-2*patrolZoneH-4, true)
	r.renderOtherZone(r.enemy, 14, 1, false)
}

func (r *tcellRenderer) renderPlayerField() {
	r.currUiLine = r.h/2 - 2
	cw.SetStyle(tcell.ColorRed, tcell.ColorBlack)
	r.drawLineAndIncrementY(fmt.Sprintf("Base HP %d", r.activePlayer.baseHealth), 0)
	cw.ResetStyle()
	for _, tb := range r.activePlayer.techBuildings {
		if tb != nil {
			line := tb.static.name
			if tb.isUnderConstruction {
				cw.SetFg(tcell.ColorRed)
				line += " (builds)"
			} else {
				cw.ResetStyle()
			}
			r.drawLineAndIncrementY(line, 0)
		}
	}
	cw.ResetStyle()
	if r.activePlayer.addonBuilding != nil {
		r.drawLineAndIncrementY(r.activePlayer.addonBuilding.static.name, 0)
	}
	r.drawLineAndIncrementY("(B)build", 0)
	r.drawLineAndIncrementY(fmt.Sprintf("DRAW: %4d", len(r.activePlayer.draw)), 0)
	r.drawLineAndIncrementY(fmt.Sprintf("DISCARD: %d", len(r.activePlayer.discard)), 0)
	if r.activePlayer.hiredWorkerThisTurn {
		cw.SetFg(tcell.ColorGreen)
	} else {
		cw.SetFg(tcell.ColorRed)
	}
	r.drawLineAndIncrementY(fmt.Sprintf("WORKERS: %d", r.activePlayer.workers), 0)
	cw.SetFg(tcell.ColorYellow)
	r.drawLineAndIncrementY(fmt.Sprintf("$%d", r.activePlayer.gold), 0)
	cw.SetStyle(tcell.ColorBlack, tcell.ColorYellow)
	r.drawLineAndIncrementY("C - access command zone", 0)
	cw.ResetStyle()
	r.renderPatrolZone(r.activePlayer, r.h-cardShortH-patrolZoneH-2, false)
	r.renderOtherZone(r.activePlayer, 14, r.h/2-2, true)
	r.renderHand()
}

func (r *tcellRenderer) renderHand() {
	cards := len(r.activePlayer.hand)
	if cards == 0 {
		return
	}
	cardW := r.w / cards
	if cardW > r.w/5 {
		cardW = r.w / 5
	}
	for i, c := range r.activePlayer.hand {
		r.renderCardInHand(c, i*(cardW), r.h-cardShortH, cardW, cardShortH)
		cw.SetStyle(tcell.ColorBlack, tcell.ColorDarkGray)
		cw.PutStringCenteredAt(strconv.Itoa(i+1), i*(cardW)+cardW/2, r.h-cardShortH)
	}
}

func (r *tcellRenderer) renderSelectedCardFromHand() {
	r.renderCardFull(r.pc.currentSelectedCardFromHand, r.w/2-cardFullW/2, r.h/2-cardFullH/2, cardFullW, cardFullH)
	cw.SetStyle(tcell.ColorBlack, tcell.ColorBlue)
	r.currUiLine = r.h/2 + cardFullH/2 + 1
	r.drawLineAndIncrementY(fmt.Sprintf(" %-30s", "W - spend as worker"), r.w/2-cardFullW/2)
	r.drawLineAndIncrementY(fmt.Sprintf(" %-30s", "P - play card"), r.w/2-cardFullW/2)
}

func (r *tcellRenderer) renderSelectedUnit() {
	r.renderCardFull(r.pc.currentSelectedUnit.card, r.w/2-cardFullW/2, r.h/2-cardFullH/2, cardFullW, cardFullH)
	r.currUiLine = r.h/2 + cardFullH/2 + 1
	cw.SetStyle(tcell.ColorBlack, tcell.ColorBlue)
	r.drawLineAndIncrementY(fmt.Sprintf(" %-30s", "M - move to other zone"), r.w/2-cardFullW/2)
	if r.g.canUnitAttack(r.pc.currentSelectedUnit) {
		r.drawLineAndIncrementY(fmt.Sprintf(" %-30s", "A - attack"), r.w/2-cardFullW/2)
	}
	// r.drawLineAndIncrementY(fmt.Sprintf(" %-30s", "U - use ability"), r.w/2-cardFullW/2)
	if r.pc.currentSelectedUnit.isHero() {
		r.drawLineAndIncrementY(fmt.Sprintf(" %-30s", "L - level up"), r.w/2-cardFullW/2)
	}
}

func (r *tcellRenderer) renderOtherZone(p *player, x, y int, renderSelectionStrings bool) {
	cw.ResetStyle()
	const keys = "QWERT     "
	for i, unt := range p.otherZone {
		str := "    "
		a, d := unt.getAtkHpWithWounds()
		if renderSelectionStrings {
			str = string(keys[i]) + " - "
		}
		str += fmt.Sprintf("%d/%d %s", a, d, unt.getName())
		cw.ResetStyle()
		if unt.tapped {
			cw.SetFg(tcell.ColorDarkBlue)
			str += " TAPPED"
		} else if unt.isHero() {
			str += fmt.Sprintf(" LVL %d", unt.level)
		}
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
			cw.SetStyle(tcell.ColorBlack, tcell.ColorYellow)
			cw.PutStringCenteredAt(fmt.Sprintf("Press %d to hire", i+1), cx+i*(cardW+2)+cardW/2, cy+cardH)
		}
	}
}

func (r *tcellRenderer) renderCodexSelection(p *player) {
	ww, wh := r.w-r.w/5, r.h-r.h/10
	wx, wy := (r.w-ww)/2, (r.h-wh)/2
	cardW, cardH := 2*ww/5, wh-2
	cardX, cardY := wx+ww-cardW, wy+1
	r.drawWindow("SELECT CARDS FROM YOUR CODEX", wx, wy, ww, wh, tcell.ColorBlue)
	if r.pc.currentSelectedCardFromHand != nil {
		r.renderCardFull(r.pc.currentSelectedCardFromHand, cardX, cardY, cardW, cardH)
	} else {
		cw.SetStyle(tcell.ColorBlack, tcell.ColorBlue)
		cw.DrawRect(cardX, cardY, 0, cardH)
		cw.ResetStyle()
		cw.PutStringCenteredAt("Select a card", cardX+cardW/2, cardY+cardH/2)
		if !p.isObligatedToAdd2Cards() {
			cw.PutStringCenteredAt("Or press ENTER to skip", cardX+cardW/2, cardY+cardH/2+1)
		}
	}
	r.currUiLine = wy + 1
	codex := p.codices[r.pc.currentCodexPage]
	cw.SetFg(tcell.ColorDarkMagenta)
	separatorText, prevSeparatorText := "", ""
	for i := range codex.cards {
		selectionLine := ""
		thisCard := codex.getCardByIndex(i)
		switch thisCard.(type) {
		case *magicCard:
			separatorText = "MAGIC"
			selectionLine = fmt.Sprintf("%s - %-25s (x%d)", string(playerCodexCardSelectionKeys[i]),
				thisCard.getName(), codex.cardsCounts[i])
		case *unitCard:
			separatorText = fmt.Sprintf("TECH %d UNITS", thisCard.(*unitCard).techLevel)
			a, d := thisCard.(*unitCard).baseAtk, thisCard.(*unitCard).baseHP
			selectionLine = fmt.Sprintf("%s - %-25s %d/%d (x%d)", string(playerCodexCardSelectionKeys[i]),
				thisCard.getName(), a, d, codex.cardsCounts[i])
		}
		// draw separator line if neccessary
		if separatorText != prevSeparatorText {
			cw.SetFg(tcell.ColorDarkMagenta)
			r.drawLineAndIncrementY(separatorText, wx+3)
		}
		prevSeparatorText = separatorText
		if codex.cardsCounts[i] > 0 {
			cw.SetFg(tcell.ColorWhite)
		} else {
			cw.SetFg(tcell.ColorDarkGray)
		}
		r.drawLineAndIncrementY(selectionLine, wx+1)
	}
}

func (r *tcellRenderer) renderPatrolZone(p *player, y int, shorten bool) {
	x := r.w - patrolZoneW - 1
	cw.ResetStyle()
	cardW := patrolZoneW / 5
	currentDrawPos := 5
	for pos := 4; pos >= 0; pos-- {
		if shorten && p.patrolZone[pos] == nil {
			continue
		}
		currentDrawPos--
		currX := x + currentDrawPos*cardW
		cw.SetStyle(tcell.ColorBlack, tcell.ColorDarkGray)
		cw.DrawRect(currX, y, cardW, patrolZoneH)
		cw.SetStyle(tcell.ColorDarkGray, tcell.ColorBlack)
		position := ""
		descrString := ""
		hotkey := ""
		switch pos {
		case 0:
			position = "Leader"
			hotkey = "Y"
			if p.patrolLeaderHasShield {
				descrString = "+SHLD/Taunt"
			} else {
				descrString = "+Taunt"
			}
		case 1:
			position = "Elite"
			hotkey = "U"
			descrString = "+1 ATK"
		case 2:
			position = "Scavenger"
			hotkey = "I"
			descrString = "Dies: +1$"
		case 3:
			position = "Technician"
			hotkey = "O"
			descrString = "Dies: +1 Card"
		case 4:
			position = "Lookout"
			hotkey = "P"
			descrString = "Resist 1"
		}
		cw.PutStringCenteredAt(position, x+currentDrawPos*cardW+cardW/2, y+patrolZoneH/2)
		cw.PutStringCenteredAt(descrString, currX+cardW/2, y+patrolZoneH)
		unitHere := p.patrolZone[pos]
		if unitHere != nil {
			r.drawUnit(unitHere, x+currentDrawPos*cardW+1, y+1, cardW-2, patrolZoneH-2)
		}
		if p == r.activePlayer {
			cw.SetStyle(tcell.ColorBlack, tcell.ColorDarkGray)
			cw.PutStringCenteredAt(hotkey, currX+cardW/2, y)
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
