package main

type game struct {
	players             [2]*player
	currentPlayer       *player
	currentPlayerNumber int
	currentTurn         int
	currentPhase        int

	messagesForPlayer string
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

	g.players[0].commandZone[0] = heroCardsDb[0]
	g.players[1].commandZone[0] = heroCardsDb[1]

	// form codices
	for _, p := range g.players {
		for i, hero := range p.commandZone {
			if hero != nil {
				for _, c := range cardsDb {
					if c.getElement() == hero.getElement() {
						// adding two of each card
						p.codices[i].addCard(c)
						p.codices[i].addCard(c)
					}
				}
			}
		}
	}
}

func (g *game) endCurrentPhase() {
	if g.currentPhase == PHASE_CODEX &&
		// don't end phase if a player has not taken their cards
		(g.currentPlayer.isObligatedToAdd2Cards() && g.currentPlayer.cardsAddedFromCodexThisTurn < 2) {
		// BUG: only one card can be added if the player is not obligated to take two.
		return
	}
	g.currentPhase++
	// end turn
	if g.currentPhase == TOTAL_PHASES {
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
	case PHASE_APPLY_TECH:
		g.applyTechPhase()
	// Phase 1: Untap
	case PHASE_READY:
		g.untapPhase()
	// Phase 2: Upkeep
	case PHASE_UPKEEP:
		g.upkeepPhase()
	// Phase 3: Main
	case PHASE_MAIN:
		// handled by player controllers by now
	// Phase 4: Discard
	case PHASE_DISCARD:
		g.discardPhase()
	// Phase 5: Select tech
	case PHASE_CODEX:
		// handled by player controllers by now
	}
}

func (g *game) applyTechPhase() {
	for i, c := range g.currentPlayer.cardsToAddNextTurn {
		if c != nil {
			g.currentPlayer.discard.addToBottom(c)
			g.currentPlayer.cardsToAddNextTurn[i] = nil
		}
	}
	g.currentPlayer.cardsAddedFromCodexThisTurn = 0
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
	g.currentPlayer.gold += g.currentPlayer.workers
	if g.currentPlayer.addonBuilding != nil && g.currentPlayer.addonBuilding.isUnderConstruction {
		g.currentPlayer.addonBuilding.isUnderConstruction = false
	}
	for _, tb := range g.currentPlayer.techBuildings {
		if tb != nil && tb.isUnderConstruction {
			tb.isUnderConstruction = false
		}
	}
}

func (g *game) discardPhase() {
	p := g.currentPlayer
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
