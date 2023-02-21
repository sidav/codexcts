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
	io.renderGame(g, 0, g.currentPhase)
	key := readKey()
	for key != "ESCAPE" {
		g.performCurrentPhase()
		io.renderGame(g, g.currentPlayerNumber, g.currentPhase)
		key = readKey()
		g.endCurrentPhase()
	}
}
