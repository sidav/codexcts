package main

import (
	"strings"
	"time"
)

const (
	PCMODE_NONE = iota
	PCMODE_CARD_FROM_HAND_SELECTED
	PCMODE_UNIT_SELECTED
	PCMODE_MOVE_SELECTED_UNIT
	PCMODE_SELECT_BUILDING
	PCMODE_SELECT_HERO_TO_PLAY

	PCMODE_CALLBACK_TARGET_SELECTION
)

var playerHandSelectionKeys = "1234567890"
var playerOtherZoneSelectionKeys = "qwert"
var playerPatrolZoneSelectionKeys = "yuiop"
var playerCodexCardSelectionKeys = "1234567890qw"

type playerController struct {
	controlsPlayer              *player
	currentMode                 int
	currentSelectedCardFromHand card

	currentSelectedUnit *unit
	selectedUnitZone    int
	selectedUnitIndex   int
	currentCodexPage    int

	callbackCoordsList []*playerZoneCoords
	callbackMessage    string
	g                  *game

	endPhase bool
}

func (pc *playerController) resetState() {
	pc.currentMode = PCMODE_NONE
	pc.currentSelectedCardFromHand = nil
	pc.currentSelectedUnit = nil
	pc.selectedUnitIndex = 0
	pc.selectedUnitZone = 0
}

func (pc *playerController) phaseEnded() bool {
	return pc.endPhase
}

func (pc *playerController) act(g *game) {
	pc.g = g
	pc.endPhase = false
	io.renderGame(g, g.currentPlayerNumber, pc)
	switch g.currentPhase {
	case PHASE_MAIN:
		pc.mainPhase(g)
	case PHASE_CODEX:
		pc.selectCardFromCodex(g)
	default:
		pc.resetState()
		pc.endPhase = true
		time.Sleep(200 * time.Millisecond)
	}
}

func (pc *playerController) mainPhase(g *game) {
	key := readKey()
	switch pc.currentMode {
	case PCMODE_NONE:
		switch key {
		case "ESCAPE":
			exitGame = true
		case "ENTER":
			pc.currentMode = PCMODE_NONE
			pc.endPhase = true
		}
		// build
		if key == "b" {
			pc.currentMode = PCMODE_SELECT_BUILDING
		}
		if key == "c" {
			pc.currentMode = PCMODE_SELECT_HERO_TO_PLAY
		}
		// number pressed (select card from hand)
		index := strings.Index(playerHandSelectionKeys, key)
		if index != -1 && index < len(pc.controlsPlayer.hand) {
			pc.currentMode = PCMODE_CARD_FROM_HAND_SELECTED
			pc.currentSelectedCardFromHand = pc.controlsPlayer.hand[index]
		}
		// pressed qwert
		index = strings.Index(playerOtherZoneSelectionKeys, key)
		if index != -1 {
			if index < len(pc.controlsPlayer.otherZone) {
				pc.currentMode = PCMODE_UNIT_SELECTED
				pc.currentSelectedUnit = pc.controlsPlayer.otherZone[index]
				pc.selectedUnitZone = PLAYERZONE_OTHER
				pc.selectedUnitIndex = index
			}
		}
		// pressed yuiop
		index = strings.Index(playerPatrolZoneSelectionKeys, key)
		if index != -1 && index < 5 && pc.controlsPlayer.patrolZone[index] != nil {
			pc.currentMode = PCMODE_UNIT_SELECTED
			pc.currentSelectedUnit = pc.controlsPlayer.patrolZone[index]
			pc.selectedUnitZone = PLAYERZONE_PATROL
			pc.selectedUnitIndex = index
		}
	case PCMODE_CARD_FROM_HAND_SELECTED:
		switch key {
		case "ESCAPE", "ENTER":
			pc.currentMode = PCMODE_NONE
		case "w":
			if g.tryPlayCardAsWorker(pc.currentSelectedCardFromHand) {
				pc.currentMode = PCMODE_NONE
				pc.currentSelectedCardFromHand = nil
			}
		case "p":
			switch pc.currentSelectedCardFromHand.(type) {
			case *unitCard:
				if g.tryPlayUnitCardFromHand(pc.currentSelectedCardFromHand) {
					pc.currentMode = PCMODE_NONE
					pc.currentSelectedCardFromHand = nil
				}
			case *magicCard:
				if g.tryPlayMagicCardFromHand(pc.currentSelectedCardFromHand) {
					pc.currentMode = PCMODE_NONE
					pc.currentSelectedCardFromHand = nil
				}
			}
		}
	case PCMODE_UNIT_SELECTED:
		switch key {
		case "ESCAPE", "ENTER":
			pc.currentMode = PCMODE_NONE
		case "m":
			pc.currentMode = PCMODE_MOVE_SELECTED_UNIT
		case "l":
			if g.tryLevelUpHero(pc.controlsPlayer, pc.currentSelectedUnit) {
				pc.currentMode = PCMODE_NONE
			}
		case "a":
			if g.tryAttackAsUnit(pc.controlsPlayer, pc.currentSelectedUnit) {
				pc.currentMode = PCMODE_NONE
			}
		}
	case PCMODE_MOVE_SELECTED_UNIT:
		switch key {
		case "ESCAPE", "ENTER":
			pc.currentMode = PCMODE_NONE
		}
		index := strings.Index(playerOtherZoneSelectionKeys, key)
		if index != -1 {
			pc.controlsPlayer.moveUnit(pc.currentSelectedUnit, pc.selectedUnitZone, pc.selectedUnitIndex, PLAYERZONE_OTHER, 0)
			pc.currentMode = PCMODE_NONE
		}
		// pressed yuiop
		index = strings.Index(playerPatrolZoneSelectionKeys, key)
		if index != -1 && index < 5 {
			pc.controlsPlayer.moveUnit(pc.currentSelectedUnit, pc.selectedUnitZone, pc.selectedUnitIndex, PLAYERZONE_PATROL, index)
			pc.currentMode = PCMODE_NONE
		}
	case PCMODE_SELECT_BUILDING:
		switch key {
		case "ESCAPE", "ENTER":
			pc.currentMode = PCMODE_NONE
		case "t": // try build next tech
			if g.tryBuildNextTechForPlayer(pc.controlsPlayer) {
				pc.currentMode = PCMODE_NONE
			}
		case "o": // try build tower
			if g.tryBuildBuildingForPlayer(pc.controlsPlayer, getBuildingStaticByName("Tower")) {
				pc.currentMode = PCMODE_NONE
			}
		case "s": // try build surplus
			if g.tryBuildBuildingForPlayer(pc.controlsPlayer, getBuildingStaticByName("Surplus")) {
				pc.currentMode = PCMODE_NONE
			}
		case "l": // try build lab
		case "h": // try build hall
		}
	case PCMODE_SELECT_HERO_TO_PLAY:
		switch key {
		case "ESCAPE", "ENTER":
			pc.currentMode = PCMODE_NONE
		}
		index := strings.Index("123", key)
		if index != -1 && pc.controlsPlayer.commandZone[index] != nil {
			if g.tryPlayHeroCard(pc.controlsPlayer.commandZone[index]) {
				pc.currentMode = PCMODE_NONE
			}
		}
	}
}

func (pc *playerController) selectCardFromCodex(g *game) {
	key := readKey()
	switch key {
	case "ESCAPE":
		exitGame = true
		return
	case "ENTER":
		if pc.currentSelectedCardFromHand != nil {
			if g.tryAddCardFromCodex(pc.controlsPlayer, pc.currentSelectedCardFromHand, pc.currentCodexPage) {
				pc.currentSelectedCardFromHand = nil
			}
		}
		pc.endPhase = true
		return
	case "r": // select random card; needed for debug
		index := rnd.Rand(len(pc.controlsPlayer.codices[pc.currentCodexPage].cards))
		for pc.controlsPlayer.codices[pc.currentCodexPage].cardsCounts[index] == 0 {
			index = rnd.Rand(len(pc.controlsPlayer.codices[pc.currentCodexPage].cards))
		}
		g.tryAddCardFromCodex(pc.controlsPlayer, pc.controlsPlayer.codices[pc.currentCodexPage].getCardByIndex(index), pc.currentCodexPage)
		pc.currentSelectedCardFromHand = nil
		pc.endPhase = true
		return
	case "LEFT":
		pc.currentSelectedCardFromHand = nil
		pc.currentCodexPage--
		if pc.currentCodexPage == -1 {
			pc.currentCodexPage = 2
		}
	case "RIGHT":
		pc.currentSelectedCardFromHand = nil
		pc.currentCodexPage++
		if pc.currentCodexPage == 3 {
			pc.currentCodexPage = 0
		}
	}

	index := strings.Index(playerCodexCardSelectionKeys, key)
	if index != -1 {
		pc.currentSelectedCardFromHand = pc.controlsPlayer.codices[pc.currentCodexPage].getCardByIndex(index)
	} else {
		pc.currentSelectedCardFromHand = nil
	}

	for pc.controlsPlayer.codices[pc.currentCodexPage].getTotalCardsCount() == 0 {
		pc.currentCodexPage--
	}
}
