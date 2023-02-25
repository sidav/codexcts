package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)

func (r *tcellRenderer) renderCardInHand(c card, x, y, w, h int) {
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
		for i := range cc.passiveAbilities {
			cw.PutStringCenteredAt(cc.passiveAbilities[i].getFormattedName(), x+w/2, y+3+i)
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
	separatorRelY := h / 4
	cw.SetStyle(tcell.ColorGray, tcell.ColorDarkGray)
	cw.DrawRect(x, y, w, h)
	cw.DrawRect(x, y, w, separatorRelY)
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
		spellDesc := ""
		if mc.isMinor {
			spellDesc = "Minor " + spellDesc
		}
		if mc.isUltimate {
			spellDesc = "Ultimate " + spellDesc
		}
		if mc.isOngoing {
			spellDesc = "Ongoing " + spellDesc
		}
		spellDesc += "Spell - " + mc.getSubtype()
		cw.PutStringCenteredAt(spellDesc, x+w/2, y+separatorRelY-1)
		cw.PutTextInRect(mc.description, x+1, y+h/4+3, w-2)
		cw.ResetStyle()
	case *unitCard:
		cc := c.(*unitCard)
		cw.SetFg(tcell.ColorGray)
		cw.PutStringCenteredAt("Unit - "+cc.getSubtype(), x+w/2, y+separatorRelY-1)
		for i := range cc.passiveAbilities {
			cw.PutStringCenteredAt(cc.passiveAbilities[i].getFormattedName(), x+w/2, y+separatorRelY+3+i)
		}
		cw.ResetStyle()
		elementAndTechLine += fmt.Sprintf(" Tech %d", cc.techLevel)
		cw.PutStringPaddedToRight(fmt.Sprintf("ATK %d HP %d", cc.baseAtk, cc.baseHP), x+w, y+h-1)
	case *heroCard:
		elementAndTechLine += " Hero"
		hc := c.(*heroCard)
		levelsStr := fmt.Sprintf("LEVELS 1-%d: %d/%d \n ", hc.levelsAttDef[1][0]-1, hc.levelsAttDef[1][1], hc.levelsAttDef[1][2])
		levelsStr += hc.levelsAbilitiesTexts[0] + " \n "
		for i := 1; i < len(hc.levelsAttDef)-1; i++ {
			levelsStr += fmt.Sprintf("LEVELS %d-%d: %d/%d \n ", hc.levelsAttDef[i][0], hc.levelsAttDef[i+1][0]-1,
				hc.levelsAttDef[i][1], hc.levelsAttDef[i][2])
			levelsStr += hc.levelsAbilitiesTexts[i] + " \n \n "
		}
		levelsStr += fmt.Sprintf("LEVEL %d: %d/%d \n ", hc.levelsAttDef[len(hc.levelsAttDef)-1][0],
			hc.levelsAttDef[len(hc.levelsAttDef)-1][1], hc.levelsAttDef[len(hc.levelsAttDef)-1][2])
		levelsStr += hc.levelsAbilitiesTexts[len(hc.levelsAttDef)-1] + " \n "
		cw.PutTextInRect(levelsStr, x+1, y+h/4+1, w-2)
	}
	cw.PutString(elementAndTechLine, x+1, y+h-1)
}

func (r *tcellRenderer) drawUnit(u *unit, x, y, w, h int) {
	if u.tapped {
		cw.SetStyle(tcell.ColorBlack, tcell.ColorDarkBlue)
	} else {
		cw.SetStyle(tcell.ColorBlack, tcell.ColorDarkGray)
	}
	cw.DrawFilledRect(' ', x, y, w, h)
	cw.PutTextInRect(u.card.getName(), x+1, y, w-2)
	atk, hp := u.getAtkHp()
	if u.isHero() {
		cw.PutStringCenteredAt(fmt.Sprintf("LEVEL %d", u.level), x+w/2, y+h-2)
	}
	if u.tapped {
		cw.PutStringCenteredAt("TAPPED", x+w/2, y+h-2)
	}
	cw.PutStringCenteredAt(fmt.Sprintf("%d/%d", atk, hp), x+w/2, y+h-1)
}
