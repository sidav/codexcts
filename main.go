package main

import (
	"codexcts/lib/random"
	"codexcts/lib/random/pcgrandom"
)

var (
	rnd random.PRNG
	pc  *playerController
)

func main() {
	onInit()
	defer onClose()

	rnd = pcgrandom.New(-1)

	g := &game{}
	g.initGame()
	pc = &playerController{
		controlsPlayer:              g.players[0],
		currentMode:                 PCMODE_NONE,
		currentSelectedCardFromHand: nil,
		currentSelectedUnit:         nil,
	}
	gameLoop(g)
}
