package main

import (
	"codexcts/lib/random"
	"codexcts/lib/random/pcgrandom"
)

var (
	rnd random.PRNG
)

func main() {
	rnd = pcgrandom.New(-1)
	onInit()
	defer onClose()

	g := &game{}
	g.initGame()
	for g.currentPhase != 3 {
		g.performCurrentPhase()
		g.endCurrentPhase()
	}
	io.renderGame(g, 0)
	readKey()
}
