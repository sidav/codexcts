package main

import "log"

type aiPlayerController struct {
	controlsPlayer *player
}

func (ai *aiPlayerController) act(g *game) {
	switch g.currentPhase {
	case PHASE_MAIN:
		ai.logTurnBeginning(g)
		ai.actMain(g)
	case PHASE_CODEX:
		ai.actCodex(g)
	}
}

func (ai *aiPlayerController) actMain(g *game) {
	plr := ai.controlsPlayer
	// first, play random card from hand as a worker
	if plr.gold > 0 && plr.hand.size() > 0 && plr.workers < 10 {
		worker := plr.hand[rnd.Rand(plr.hand.size())]
		log.Printf("I played %s as worker.\n", worker.getName())
		g.tryPlayCardAsWorker(worker)
	}
	// second, play random creature
	if rnd.OneChanceFrom(3) {
		var cardToPlay card
		for _, c := range plr.hand {
			if cardToPlay == nil || rnd.OneChanceFrom(3) && g.canPlayerPlayCard(plr, c) {
				cardToPlay = c
			}
		}
		log.Printf("I played %s from my hand.\n", cardToPlay.getName())
		g.tryPlayUnitCardFromHand(cardToPlay)
	}
	// third, play hero
	if rnd.OneChanceFrom(4) {
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
			g.tryPlayHeroCard(cardToPlay)
		} else {
			log.Printf("  I didn't found any heroes\n")
		}
	}
	// fourth, move creatures to patrol zone
	if len(plr.otherZone) > 0 && rnd.OneChanceFrom(2) {
		indexToMove := rnd.Rand(len(plr.otherZone))
		unitToMove := plr.otherZone[indexToMove]
		for i := range plr.patrolZone {
			if plr.patrolZone[i] == nil && rnd.OneChanceFrom(5) {
				plr.moveUnit(unitToMove, PLAYERZONE_OTHER, indexToMove, PLAYERZONE_PATROL, i)
				break
			}
		}
		log.Printf("I moved %s to other zone.\n", unitToMove.card.getName())
	}
	// fifth, build
	// TODO
	log.Printf("I ended my main phase. \n")
}

func (ai *aiPlayerController) actCodex(g *game) {
	plr := ai.controlsPlayer
	if plr.workers > 10 && !rnd.OneChanceFrom(10) {
		log.Printf("I'm not adding any more cards\n")
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

func (ai *aiPlayerController) logTurnBeginning(g *game) {
	log.Printf("===== AI TURN =====\n")
	log.Printf("Main phase of turn %d\n", g.currentTurn)
	log.Printf("My hand is: \n")
	for _, c := range ai.controlsPlayer.hand {
		log.Printf("   %s", c.getName())
	}
}
