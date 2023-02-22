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
	if pc.currentMode == PCMODE_CARD_FROM_HAND_SELECTED {
		r.renderCardFromHand()
	}
	if pc.currentMode == PCMODE_UNIT_SELECTED {
		r.renderSelectedUnit()
	}

	cw.FlushScreen()
}

func (r *tcellRenderer) renderHeader() {
	cw.SetStyle(tcell.ColorBlack, tcell.ColorYellow)
	cw.DrawRect(0, 0, r.w, 0)
	cw.SetStyle(tcell.ColorYellow, tcell.ColorBlack)
	cw.PutStringCenteredAt(fmt.Sprintf(" PLAYER %d - %s phase ", r.g.currentPlayerNumber, r.g.getCurrentPhaseName()), r.w/2, 0)
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
		r.renderCardShort(c, i*(cardW), r.h-cardShortH, cardW, cardShortH)
	}
}

func (r *tcellRenderer) renderCardShort(c card, x, y, w, h int) {
	cw.SetStyle(tcell.ColorGray, tcell.ColorDarkGray)
	cw.DrawRect(x, y, w, h)
	cw.SetStyle(tcell.ColorBlack, tcell.ColorYellow)
	cw.PutString(fmt.Sprintf("$%d", c.getCost()), x+1, y+1)
	cw.ResetStyle()
	cw.PutTextInRect(" "+c.getName(), x+3, y+1, w-6)
	elementAndTechLine := c.getElement().getName()
	switch c.(type) {
	case *magicCard:
		mc := c.(*magicCard)
		elementAndTechLine += " Magic"
		cw.SetFg(tcell.ColorGray)
		cw.PutTextInRect(mc.description, x+1, y+2, w-2)
		cw.ResetStyle()
	case *unitCard:
		cc := c.(*unitCard)
		cw.SetFg(tcell.ColorGray)
		for i := range cc.specials {
			cw.PutStringCenteredAt(cc.specials[i].getFormattedName(), x+w/2, y+3+i)
		}
		cw.ResetStyle()
		elementAndTechLine += fmt.Sprintf(" Tech %d", cc.techLevel)
		cw.PutStringPaddedToRight(fmt.Sprintf("%d/%d", cc.baseAtk, cc.baseHP), x+w, y+h-2)
	}
	cw.PutString(elementAndTechLine, x+1, y+h-1)
}

func (r *tcellRenderer) renderCardFull(c card, x, y, w, h int) {
	cw.ResetStyle()
	cw.DrawFilledRect(' ', x, y, w, h)
	cw.SetStyle(tcell.ColorGray, tcell.ColorDarkGray)
	cw.DrawRect(x, y, w, h)
	cw.DrawRect(x, y, w, h/4)
	cw.SetStyle(tcell.ColorBlack, tcell.ColorYellow)
	cw.PutString(fmt.Sprintf(" $%d ", c.getCost()), x+1, y+1)
	cw.ResetStyle()
	cw.PutTextInRect(" "+c.getName(), x+3, y+2, w-6)
	elementAndTechLine := c.getElement().getName()
	switch c.(type) {
	case *magicCard:
		mc := c.(*magicCard)
		elementAndTechLine += " Magic"
		cw.SetFg(tcell.ColorGray)
		cw.PutTextInRect(mc.description, x+1, y+h/4+3, w-2)
		cw.ResetStyle()
	case *unitCard:
		cc := c.(*unitCard)
		cw.SetFg(tcell.ColorGray)
		for i := range cc.specials {
			cw.PutStringCenteredAt(cc.specials[i].getFormattedName(), x+w/2, y+h/4+3+i)
		}
		cw.ResetStyle()
		elementAndTechLine += fmt.Sprintf(" Tech %d", cc.techLevel)
		cw.PutStringPaddedToRight(fmt.Sprintf("ATK %d HP %d", cc.baseAtk, cc.baseHP), x+w, y+h-1)
	}
	cw.PutString(elementAndTechLine, x+1, y+h-1)
}

func (r *tcellRenderer) renderCardFromHand() {
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
		if p == r.activePlayer {
			cw.PutStringCenteredAt(hotkey, currX+cardW/2, y+patrolZoneH/2)
		}
	}
}

func (r *tcellRenderer) drawLineAndIncrementY(line string, x int) {
	cw.PutString(line, x, r.currUiLine)
	r.currUiLine++
}
