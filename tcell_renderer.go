package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)

const (
	patrolZoneW = 70
	patrolZoneH = 10
	cardInHandH = 8
)

type tcellRenderer struct {
	w, h         int
	g            *game
	activePlayer *player
	enemy        *player
}

func (r *tcellRenderer) renderGame(g *game, renderForPlayerNum int) {
	r.g = g
	r.activePlayer = g.players[renderForPlayerNum]
	r.enemy = g.players[(renderForPlayerNum+1)%2]
	r.w, r.h = cw.GetConsoleSize()
	cw.ClearScreen()

	r.renderEnemyField()
	r.renderPlayerField()

	cw.FlushScreen()
}

func (r *tcellRenderer) renderEnemyField() {
	cw.ResetStyle()
	cw.PutStringPaddedToRight(fmt.Sprintf("DRAW: %d", len(r.enemy.draw)), r.w, 0)
	cw.PutStringPaddedToRight(fmt.Sprintf("DISCARD: %d", len(r.enemy.discard)), r.w, 1)
	cw.PutStringCenteredAt(fmt.Sprintf("HAND: %d", len(r.enemy.hand)), r.w/2, 0)
	cw.PutString(fmt.Sprintf("WORKERS: %d", r.enemy.workers), 0, 0)
	cw.SetFg(tcell.ColorYellow)
	cw.PutString(fmt.Sprintf("$%d", r.enemy.gold), 0, 1)
	r.renderPatrolZone(r.enemy, 4)
}

func (r *tcellRenderer) renderPlayerField() {
	cw.ResetStyle()
	cw.PutStringPaddedToRight(fmt.Sprintf("DRAW: %d", len(r.activePlayer.draw)), r.w, r.h/2)
	cw.PutStringPaddedToRight(fmt.Sprintf("DISCARD: %d", len(r.activePlayer.discard)), r.w, r.h/2+1)
	cw.PutString(fmt.Sprintf("WORKERS: %d", r.activePlayer.workers), 0, r.h/2+1)
	cw.SetFg(tcell.ColorYellow)
	cw.PutString(fmt.Sprintf("$%d", r.activePlayer.gold), 0, r.h/2)
	r.renderPatrolZone(r.activePlayer, r.h-cardInHandH-patrolZoneH-2)
	r.renderHand()
}

func (r *tcellRenderer) renderHand() {
	cards := len(r.activePlayer.hand)
	cardW := r.w / cards
	if cardW > r.w/5 {
		cardW = r.w / 5
	}
	for i, c := range r.activePlayer.hand {
		r.renderCardShort(c, i*(cardW), r.h-cardInHandH, cardW, cardInHandH)
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
		cw.PutStringPaddedToRight(fmt.Sprintf("%d/%d", cc.baseAtk, cc.baseDef), x+w, y+h-1)
	}
	cw.PutString(elementAndTechLine, x+1, y+h-2)
}

func (r *tcellRenderer) renderPatrolZone(p *player, y int) {
	x := r.w - patrolZoneW - 12
	cw.ResetStyle()
	cardW := patrolZoneW / 5
	for i := 0; i < 5; i++ {
		currX := x + i*cardW
		cw.SetStyle(tcell.ColorBlack, tcell.ColorGray)
		cw.DrawRect(currX, y, cardW, patrolZoneH)
		cw.SetStyle(tcell.ColorGray, tcell.ColorBlack)
		switch i {
		case 0:
			cw.PutStringCenteredAt("Squad leader", currX+cardW/2, y+1)
			cw.PutStringCenteredAt("+SHLD/Taunt", currX+cardW/2, y+patrolZoneH)
		case 1:
			cw.PutStringCenteredAt("Elite", currX+cardW/2, y+1)
			cw.PutStringCenteredAt("+1 ATK", currX+cardW/2, y+patrolZoneH)
		case 2:
			cw.PutStringCenteredAt("Scavenger", x+i*cardW+cardW/2, y+1)
			cw.PutStringCenteredAt("Dies: +1$", currX+cardW/2, y+patrolZoneH)
		case 3:
			cw.PutStringCenteredAt("Technician", x+i*cardW+cardW/2, y+1)
			cw.PutStringCenteredAt("Dies: +1 Card", currX+cardW/2, y+patrolZoneH)
		case 4:
			cw.PutStringCenteredAt("Lookout", x+i*cardW+cardW/2, y+1)
			cw.PutStringCenteredAt("Resist 1", currX+cardW/2, y+patrolZoneH)
		}
	}
}
