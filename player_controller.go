package main

import (
	"strconv"
	"strings"
	"time"
)

const (
	PCMODE_NONE = iota
	PCMODE_CARD_FROM_HAND_SELECTED
	PCMODE_UNIT_SELECTED
	PCMODE_MOVE_SELECTED_UNIT
)

var playerOtherZoneSelectionKeys = "qwert"
var playerPatrolZoneSelectionKeys = "yuiop"

type playerController struct {
	controlsPlayer              *player
	currentMode                 int
	currentSelectedCardFromHand card

	currentSelectedUnit *unit
	selectedUnitZone    int
	selectedUnitIndex   int

	exitGame   bool
	phaseEnded bool
}

func (pc *playerController) resetState() {
	pc.currentMode = PCMODE_NONE
	pc.currentSelectedCardFromHand = nil
	pc.currentSelectedUnit = nil
	pc.selectedUnitIndex = 0
	pc.selectedUnitZone = 0
}

func (pc *playerController) act(g *game) {
	switch g.currentPhase {
	case 3:
		pc.mainPhase(g)
	default:
		time.Sleep(250 * time.Millisecond)
		pc.phaseEnded = true
	}
}

func (pc *playerController) mainPhase(g *game) {
	key := readKey()
	switch pc.currentMode {
	case PCMODE_NONE:
		switch key {
		case "ESCAPE":
			pc.exitGame = true
		case "ENTER":
			pc.currentMode = PCMODE_NONE
			pc.phaseEnded = true
		}
		// number pressed
		if index, err := strconv.Atoi(key); err == nil && len(key) == 1 {
			index--
			if index < len(pc.controlsPlayer.hand) {
				pc.currentMode = PCMODE_CARD_FROM_HAND_SELECTED
				pc.currentSelectedCardFromHand = pc.controlsPlayer.hand[index]
			}
		}
		// pressed qwert
		index := strings.Index(playerOtherZoneSelectionKeys, key)
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
			}
		}
	case PCMODE_UNIT_SELECTED:
		switch key {
		case "ESCAPE", "ENTER":
			pc.currentMode = PCMODE_NONE
		case "m":
			pc.currentMode = PCMODE_MOVE_SELECTED_UNIT
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
	}
}
