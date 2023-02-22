package main

import (
	"strconv"
	"strings"
)

const (
	PCMODE_NONE = iota
	PCMODE_CARD_FROM_HAND_SELECTED
	PCMODE_UNIT_SELECTED
)

var playerOtherZoneSelectionKeys = "qwert"
var playerPatrolZoneSelectionKeys = "yuiop"

type playerController struct {
	controlsPlayer              *player
	currentMode                 int
	currentSelectedCardFromHand card
	currentSelectedUnit         *unit
	exitGame                    bool
	phaseEnded                  bool
}

func (pc *playerController) act(g *game) {
	switch g.currentPhase {
	case 3:
		pc.mainPhase(g)
	default:
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
			}
		}
		// pressed yuiop
		index = strings.Index(playerPatrolZoneSelectionKeys, key)
		if index != -1 && index < 5 && pc.controlsPlayer.patrolZone[index] != nil {
			pc.currentMode = PCMODE_UNIT_SELECTED
			pc.currentSelectedUnit = pc.controlsPlayer.patrolZone[index]
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
		}
	}
}
