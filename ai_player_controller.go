package main

import "log"

type aiPlayerController struct {
	controlsPlayer *player
}

func (ai *aiPlayerController) act(g *game) {
	switch g.currentPhase {
	case PHASE_MAIN:
		log.Printf("===== AI TURN %d =====\n", g.currentTurn)
		log.Printf("Main phase\n")
		ai.logHand()
		ai.actMain(g)
	case PHASE_CODEX:
		ai.actCodex(g)
	}
}

func (ai *aiPlayerController) actMain(g *game) {
	// plr := ai.controlsPlayer
	// first, play random card from hand as a worker
	ai.addWorker(g)
	performedActions := 0
	for performedActions < 3 {
		if ai.tryPerformRandomAction(g) {
			performedActions++
		}
	}
	ai.tryPerformRandomAction(g)
	log.Printf("I ended my main phase with $%d. \n", ai.controlsPlayer.gold)
}

func (ai *aiPlayerController) addWorker(g *game) {
	plr := ai.controlsPlayer
	if plr.gold > 0 && plr.hand.size() > 0 && plr.workers < 10 {
		worker := plr.hand[rnd.Rand(plr.hand.size())]
		log.Printf("I played %s as worker.\n", worker.getName())
		g.tryPlayCardAsWorker(worker)
	}
}

func (ai *aiPlayerController) tryPerformRandomAction(g *game) bool {
	// TODO: build
	actionToPerform := rnd.Rand(5)
	switch {
	case actionToPerform < 1:
		return ai.tryPlayUnit(g)
	case actionToPerform < 2:
		return ai.playHero(g)
	case actionToPerform < 3:
		return ai.tryMoveUnit(g)
	case actionToPerform < 4:
		return ai.tryBuild(g)
	default:
		log.Println("I skipped an action.")
		return true
	}
}

func (ai *aiPlayerController) tryPlayUnit(g *game) bool {
	ai.logHand()
	plr := ai.controlsPlayer
	var cardToPlay card
	for _, c := range plr.hand {
		if g.canPlayerPlayCard(plr, c) && (cardToPlay == nil || rnd.OneChanceFrom(3)) {
			cardToPlay = c
		}
	}
	if cardToPlay != nil {
		log.Printf("I played %s from my hand.\n", cardToPlay.getName())
		return g.tryPlayUnitCardFromHand(cardToPlay)
	} else {
		log.Printf("I found no cards from my hand to play.\n")
		return false
	}
}

func (ai *aiPlayerController) playHero(g *game) bool {
	plr := ai.controlsPlayer
	log.Printf("I try to play hero...\n")
	var cardToPlay *heroCard
	for _, c := range plr.commandZone {
		if c != nil && g.canPlayerPlayCard(plr, c) && (cardToPlay == nil || rnd.OneChanceFrom(3)) {
			log.Printf("  I selected %s as my hero candidate.\n", c.getName())
			cardToPlay = c
		}
	}
	if cardToPlay != nil {
		log.Printf("  I played %s as my hero.\n", cardToPlay.getName())
		return g.tryPlayHeroCard(cardToPlay)
	} else {
		log.Printf("  I didn't found any heroes\n")
		return false
	}
}

func (ai *aiPlayerController) tryMoveUnit(g *game) bool {
	plr := ai.controlsPlayer
	if len(plr.otherZone) == 0 {
		log.Printf("There are no units to move.\n")
		return false
	}
	indexToMove := rnd.Rand(len(plr.otherZone))
	unitToMove := plr.otherZone[indexToMove]
	for i := range plr.patrolZone {
		if plr.patrolZone[i] == nil && rnd.OneChanceFrom(5) {
			plr.moveUnit(unitToMove, PLAYERZONE_OTHER, indexToMove, PLAYERZONE_PATROL, i)
			break
		}
	}
	log.Printf("I moved %s to other zone.\n", unitToMove.card.getName())
	return true
}

func (ai *aiPlayerController) tryBuild(g *game) bool {
	randBuilding := sTableBuildings[rnd.Rand(len(sTableBuildings))-1]
	if g.canPlayerBuild(ai.controlsPlayer, randBuilding) {
		g.tryBuildBuildingForPlayer(ai.controlsPlayer, randBuilding)
		log.Printf("I built %s.\n", randBuilding.name)
		return true
	}
	log.Println("I can't build anything.")
	return false
}

func (ai *aiPlayerController) actCodex(g *game) {
	plr := ai.controlsPlayer
	if plr.workers > 10 && !rnd.OneChanceFrom(10) {
		log.Printf("I'm not adding any more cards.\n")
		return
	}
	codexIndex := 99
	for codexIndex > 2 || plr.codices[codexIndex].getTotalCardsCount() == 0 {
		codexIndex = rnd.Rand(3)
	}
	indexOfCard := 99
	for indexOfCard > plr.codices[codexIndex].getUniqueCardsCount() {
		indexOfCard = rnd.Rand(plr.codices[codexIndex].getUniqueCardsCount())
	}
	cardToAdd := plr.codices[codexIndex].getCardByIndex(indexOfCard)
	log.Printf("I add %s from my codex\n", cardToAdd.getName())
	g.tryAddCardFromCodex(plr, cardToAdd, codexIndex)
}

func (ai *aiPlayerController) logHand() {
	log.Printf("I have $%d, My hand is: \n", ai.controlsPlayer.gold)
	for _, c := range ai.controlsPlayer.hand {
		log.Printf("   %s", c.getName())
	}
}
