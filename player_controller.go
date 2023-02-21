package main

import "strconv"

const (
	PCMODE_NONE = iota
	PCMODE_CARD_FROM_HAND_SELECTED
	PCMODE_UNIT_SELECTED
)

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
		case "1", "2", "3", "4", "5", "6":
			index, err := strconv.Atoi(key)
			if err != nil {
				return
			}
			index--
			if index < len(pc.controlsPlayer.hand) {
				pc.currentMode = PCMODE_CARD_FROM_HAND_SELECTED
				pc.currentSelectedCardFromHand = pc.controlsPlayer.hand[index]
			}
		}
	case PCMODE_CARD_FROM_HAND_SELECTED:
		switch key {
		case "ESCAPE":
			pc.exitGame = true
		case "ENTER":
			pc.currentMode = PCMODE_NONE
		case "w":
			if g.tryPlayCardAsWorker(pc.currentSelectedCardFromHand) {
				pc.currentMode = PCMODE_NONE
				pc.currentSelectedCardFromHand = nil
			}
		}
	}
}
