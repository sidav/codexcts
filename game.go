package main

type game struct {
	players [2]*player
}

func (g *game) initGame() {
	g.players[0] = &player{}
	g.players[1] = &player{}

	for i := range g.players {
		// fmt.Println("Player ", i)
		g.players[i].workers = 4
		g.players[i].baseHealth = 20

		baseDeck := getStartingCardsForElement(ELEMENT_NEUTRAL)
		// fmt.Println("Got ", len(baseDeck), " cards")
		for _, c := range baseDeck {
			g.players[i].draw.addToBottom(c)
		}
		// fmt.Printf("Deck: %d cards\n", len(g.players[i].draw))
		g.players[i].draw.shuffle(rnd)
		// fmt.Printf("Shuffled deck: %d cards\n", len(g.players[i].draw))
		for n := 0; n < 5; n++ {
			g.players[i].drawCard()
		}
		g.players[i].sortHand()
	}
	g.players[1].workers = 5
}

func (g *game) makeTurnAsPlayer(p *player) {
	// Phase 0: Apply tech
	// Phase 1: Untap
	// Phase 2: Upkeep
	// Phase 3: Main
	// Phase 4: Discard
	g.discardPhase(p)
	// Phase 5: Select tech
}

func (g *game) discardPhase(p *player) {
	cardsToDraw := len(p.hand) + 2
	if cardsToDraw >= 5 {
		cardsToDraw = 5
	}
	p.discardHand()
	for i := 0; i < cardsToDraw; i++ {
		if len(p.draw) == 0 {
			if len(p.discard) == 0 {
				break
			}
			p.addDiscardIntoDraw()
			p.shuffleDraw()
		}
		p.drawCard()
	}
}
