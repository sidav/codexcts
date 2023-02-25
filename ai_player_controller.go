package main

import "log"

type aiPlayerController struct {
	controlsPlayer *player
}

func (ai *aiPlayerController) phaseEnded() bool {
	return true
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
	const ACTIONS_PER_TURN = 5
	// first, play random card from hand as a worker
	ai.addWorker(g)
	performedActions := 0
	for performedActions < ACTIONS_PER_TURN {
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
	actionToPerform := rnd.Rand(10)
	switch {
	case actionToPerform == 0:
		return ai.tryPlayUnit(g)
	case actionToPerform == 1:
		return ai.playHero(g)
	case actionToPerform == 2:
		return ai.tryMoveUnit(g)
	case actionToPerform == 3:
		return ai.tryBuild(g)
	case actionToPerform == 4:
		return ai.tryLevelUpHero(g)
	case actionToPerform < 7:
		return ai.tryAttack(g)
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
	randBuilding := sTableBuildings[rnd.Rand(len(sTableBuildings))]
	if g.canPlayerBuild(ai.controlsPlayer, randBuilding) {
		log.Printf("I try to build %s.\n", randBuilding.name)
		g.tryBuildBuildingForPlayer(ai.controlsPlayer, randBuilding)
		log.Printf("I built %s.\n", randBuilding.name)
		return true
	}
	log.Println("I can't build anything.")
	return false
}

func (ai *aiPlayerController) tryLevelUpHero(g *game) bool {
	var heroCandidates []*unit
	for _, u := range ai.controlsPlayer.otherZone {
		if u.isHero() {
			heroCandidates = append(heroCandidates, u)
		}
	}
	for _, u := range ai.controlsPlayer.patrolZone {
		if u != nil && u.isHero() {
			heroCandidates = append(heroCandidates, u)
		}
	}
	if len(heroCandidates) > 0 {
		hero := heroCandidates[rnd.Rand(len(heroCandidates))]
		log.Printf("I try to level up %s.\n", hero.card.getName())
		return g.tryLevelUpHero(ai.controlsPlayer, hero)
	} else {
		log.Println("I can't level up anything.")
		return false
	}
}

func (ai *aiPlayerController) tryAttack(g *game) bool {
	var candidates []*unit
	for _, u := range ai.controlsPlayer.otherZone {
		if g.canUnitAttack(u) {
			candidates = append(candidates, u)
		}
	}
	if len(candidates) == 0 {
		log.Println("I have nothing to attack with.")
		return false
	} else {
		attacker := candidates[rnd.Rand(len(candidates))]
		log.Printf("I attack with %s.", attacker.card.getName())
		return g.tryAttackAsUnit(ai.controlsPlayer, attacker)
	}
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
	for indexOfCard > plr.codices[codexIndex].getRemainingUniqueCardsCount() {
		indexOfCard = rnd.Rand(plr.codices[codexIndex].getRemainingUniqueCardsCount())
	}
	cardToAdd := plr.codices[codexIndex].getCardByIndex(indexOfCard)
	log.Printf("I add %s from my codex\n", cardToAdd.getName())
	g.tryAddCardFromCodex(plr, cardToAdd, codexIndex)
}

func (ai *aiPlayerController) logHand() {
	log.Printf("I have $%d and %d workers. My hand is: \n", ai.controlsPlayer.gold, ai.controlsPlayer.workers)
	for _, c := range ai.controlsPlayer.hand {
		log.Printf("   %s", c.getName())
	}
}
