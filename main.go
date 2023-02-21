package main

import (
	"codexcts/lib/random"
	"codexcts/lib/random/pcgrandom"
	"fmt"
)

var rnd random.PRNG

func main() {
	rnd = pcgrandom.New(-1)
	g := game{}
	g.initGame()
	for pn, p := range g.players {
		fmt.Printf("Player %d:\n", pn+1)
		fmt.Printf("Hand: \n")
		for i := range p.hand {
			fmt.Printf("  %s\n", p.hand[i].getFormattedName())
		}
	}
}
