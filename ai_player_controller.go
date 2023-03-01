package main

import "log"

type aiPlayerController struct {
	controlsPlayer *player
	isItsTurn      bool
}

func (ai *aiPlayerController) phaseEnded() bool {
	return true
}

func (ai *aiPlayerController) act(g *game) {
	switch g.currentPhase {
	case PHASE_MAIN:
		ai.isItsTurn = true
		log.Println("")
		log.Printf("======= AI: %s TURN %d =======\n", ai.controlsPlayer.name, g.currentTurn)
		log.Printf("Main phase\n")
		ai.logFullReport()
		ai.actMain(g)
	case PHASE_CODEX:
		ai.actCodex(g)
		ai.isItsTurn = false

		// cheats for AI. I couldn't make it play better yet. :(
		if rnd.OneChanceFrom(3) {
			ai.controlsPlayer.gold++
		}
		if rnd.OneChanceFrom(2) && ai.controlsPlayer.hand.size() < 5 {
			ai.controlsPlayer.drawCard()
		}
	}
}

func (ai *aiPlayerController) actMain(g *game) {
	plr := ai.controlsPlayer
	const economicActionsPerTurn = 7
	// first, play random card from hand as a worker
	ai.tryAddWorker(g)
	performedEconomicActions := 0
	for performedEconomicActions < economicActionsPerTurn {
		if ai.tryPerformRandomEconomicAction(g) {
			performedEconomicActions++
		}
	}
	// attack-related actions
	attackActions := 10 // 3 + rnd.Rand(5)
	for i := 0; i < attackActions; i++ {
		ai.tryAttack(g)
	}
	ai.tryMoveUnits(g)
	ai.tryMoveUnits(g)

	log.Printf("I ended my main phase with $%d. \n", plr.gold)
	log.Printf("Hand size: %d, draw size: %d, discard size: %d \n",
		plr.hand.size(), plr.draw.size(), plr.discard.size())
}

func (ai *aiPlayerController) tryAddWorker(g *game) bool {
	plr := ai.controlsPlayer
	if plr.gold == 0 {
		return false
	}

	if plr.workers >= 10 && rnd.Rand(plr.workers*3) > 0 {
		return false
	}

	index := rnd.SelectRandomIndexFromWeighted(len(plr.hand), func(i int) int {
		switch plr.hand[i].(type) {
		case *magicCard:
			if g.canPlayerPlayCard(plr, plr.hand[i]) {
				return 5
			} else {
				return 10 // Unplayable by AI == unimplemented?
			}
		case *unitCard:
			if g.canPlayerPlayCard(plr, plr.hand[i]) {
				return 5 // playable cards are non-tech-dependent, so not important?
			} else {
				return 1
			}
		}
		return 0
	})
	worker := plr.hand[index]
	if g.tryPlayCardAsWorker(worker) {
		log.Printf("I played %s as %dth worker.\n", worker.getName(), plr.workers)
		return true
	} else {
		log.Printf("I can't play %s as worker! \n", worker.getName())
		return false
	}
}

func (ai *aiPlayerController) tryPlayMagicCard(g *game) bool {
	plr := ai.controlsPlayer
	var cardToPlay card
	for _, c := range plr.hand {
		if _, isMagic := c.(*magicCard); isMagic {
			cardToPlay = c
		}
	}
	if cardToPlay != nil {
		log.Printf("I played magic %s from my hand.\n", cardToPlay.getName())
		return g.tryPlayMagicCardFromHand(cardToPlay)
	} else {
		log.Printf("I found no magic cards from my hand to play.\n")
		return false
	}
}

func (ai *aiPlayerController) tryPlayUnit(g *game) bool {
	// ai.logHand()
	plr := ai.controlsPlayer
	var cardToPlay card
	for _, c := range plr.hand {
		if _, isMagic := c.(*magicCard); isMagic {
			continue
		}
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
		log.Printf("I try to level up %s.\n", hero.getName())
		return g.tryLevelUpHero(ai.controlsPlayer, hero)
	} else {
		log.Println("I can't level up anything.")
		return false
	}
}

func (ai *aiPlayerController) actCodex(g *game) {
	plr := ai.controlsPlayer
	if plr.workers > 10 && !rnd.OneChanceFrom(3) {
		log.Printf("I'm not adding any more cards.\n")
		return
	}
	codexIndex := rnd.SelectRandomIndexFromWeighted(3, func(x int) int {
		return plr.codices[x].getTotalCardsCount()
	})
	indexOfCard := rnd.SelectRandomIndexFromWeighted(len(plr.codices[codexIndex].cards), func(x int) int {
		if plr.codices[codexIndex].cardsCounts[x] == 0 {
			return 0
		}
		switch plr.codices[codexIndex].getCardByIndex(x).(type) {
		case *magicCard:
			return 1
		case *unitCard:
			uc := plr.codices[codexIndex].getCardByIndex(x).(*unitCard)
			weight := 1
			if plr.hasTechLevel(uc.techLevel) {
				weight += 2
			}
			if len(uc.passiveAbilities) > 0 {
				weight++
			}
			return weight
		}
		panic("Weight selection error")
	})
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

func (ai *aiPlayerController) logFullReport() {
	plr := ai.controlsPlayer
	log.Printf("Buildings: Base %d/20. ", plr.baseHealth)
	for i, b := range plr.techBuildings {
		if b != nil {
			log.Printf(" T%d: %s (%dhp)", i, b.static.name, b.currentHitpoints)
		}
	}
	if plr.addonBuilding != nil {
		log.Printf(" %s: %dhp", plr.addonBuilding.static.name, plr.addonBuilding.currentHitpoints)
	}
	log.Printf("Units in other zone:")
	for _, u := range plr.otherZone {
		log.Println(u.getNameWithStats())
	}
	log.Printf("Units in patrol zone:")
	for i, u := range plr.patrolZone {
		if u != nil {
			log.Printf("%d: %s", i, u.getNameWithStats())
		} else {
			log.Printf("%d: --", i)
		}
	}
	log.Printf("Hand size: %d, draw size: %d, discard size: %d \n",
		plr.hand.size(), plr.draw.size(), plr.discard.size())
	ai.logHand()
}
