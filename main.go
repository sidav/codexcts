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
	//for pn, p := range g.players {
	//	fmt.Printf("Player %d:\n", pn+1)
	//	fmt.Printf("Hand: \n")
	//	for i := range p.hand {
	//		fmt.Printf("  %s\n", p.hand[i].getFormattedName())
	//	}
	//}
	//g.performCurrentPhase(g.players[0])
	for g.currentPhase != 3 {
		g.performCurrentPhase()
		g.endCurrentPhase()
	}
	io.renderGame(g, 0)
	readKey()
}
