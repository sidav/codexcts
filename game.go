package main

type game struct {
	players             [2]*player
	currentPlayer       *player
	currentPlayerNumber int
	currentTurn         int
	currentPhase        int
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
	g.currentTurn = 1
	g.currentPlayer = g.players[0]
	g.currentPlayerNumber = 0
	g.currentPhase = 0
}

func (g *game) getCurrentPhaseName() string {
	switch g.currentPhase {
	case 0:
		return "Apply tech"
	case 1:
		return "Ready"
	case 2:
		return "Upkeep"
	case 3:
		return "Main"
	case 4:
		return "Discard"
	case 5:
		return "Select tech"
	}
	panic("No phase name")
}

func (g *game) endCurrentPhase() {
	g.currentPhase++
	if g.currentPhase > 5 {
		g.currentPhase = 0
		g.currentPlayerNumber = (g.currentPlayerNumber + 1) % 2
		g.currentPlayer = g.players[g.currentPlayerNumber]
		if g.currentPlayerNumber == 0 {
			g.currentTurn++
		}
	}
}

func (g *game) performCurrentPhase() {
	switch g.currentPhase {
	// Phase 0: Apply tech
	case 0:
		g.applyTechPhase()
	// Phase 1: Untap
	case 1:
		g.untapPhase()
	// Phase 2: Upkeep
	case 2:
		g.upkeepPhase()
	// Phase 3: Main
	case 3:
		// handled by player controllers by now

	// Phase 4: Discard
	case 4:
		g.discardPhase()

	// Phase 5: Select tech
	case 5:
	}
}

func (g *game) applyTechPhase() {
	for i, c := range g.currentPlayer.cardsToAddNextTurn {
		if c != nil {
			g.currentPlayer.discard.addToBottom(c)
			g.currentPlayer.cardsToAddNextTurn[i] = nil
		}
	}
}

func (g *game) untapPhase() {
	g.currentPlayer.hiredWorkerThisTurn = false
	for _, u := range g.currentPlayer.patrolZone {
		if u != nil {
			u.tapped = false
		}
	}
	for _, u := range g.currentPlayer.otherZone {
		if u != nil {
			u.tapped = false
		}
	}
}

func (g *game) upkeepPhase() {
	g.players[g.currentPlayerNumber].gold += g.players[g.currentPlayerNumber].workers
}

func (g *game) discardPhase() {
	p := g.players[g.currentPlayerNumber]
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
	p.sortHand()
}
