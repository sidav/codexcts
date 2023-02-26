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
		plr := ai.controlsPlayer
		log.Println("")
		log.Printf("======= AI: %s TURN %d =======\n", ai.controlsPlayer.name, g.currentTurn)
		log.Printf("Hand size: %d, draw size: %d, discard size: %d \n",
			plr.hand.size(), plr.draw.size(), plr.discard.size())
		log.Printf("Main phase\n")
		ai.logHand()
		ai.actMain(g)
	case PHASE_CODEX:
		ai.actCodex(g)
		ai.isItsTurn = false
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
	for i := 0; i < 5; i++ {
		actionTaken := ai.tryAttack(g)
		actionTaken = actionTaken || ai.tryMoveUnits(g)
		if !actionTaken {
			break
		}
	}
	log.Printf("I ended my main phase with $%d. \n", plr.gold)
	log.Printf("Hand size: %d, draw size: %d, discard size: %d \n",
		plr.hand.size(), plr.draw.size(), plr.discard.size())
}

func (ai *aiPlayerController) tryAddWorker(g *game) bool {
	plr := ai.controlsPlayer
	if plr.gold == 0 {
		return false
	}

	// TODO: change when the AI will be able to cast spells
	hasCardsWhichAiCantPlay := false
	for _, c := range plr.hand {
		if _, ok := c.(*magicCard); ok {
			hasCardsWhichAiCantPlay = true
			break
		}
	}

	if !hasCardsWhichAiCantPlay && plr.workers >= 10 && rnd.Rand(plr.workers*3) > 0 {
		return false
	}

	index := rnd.SelectRandomIndexFromWeighted(len(plr.hand), func(i int) int {
		switch plr.hand[i].(type) {
		case *magicCard:
			return 100 // TODO: change when the AI will be able to cast spells
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

func (ai *aiPlayerController) tryMoveUnits(g *game) bool {
	plr := ai.controlsPlayer
	moved := false
	// Maybe move something from patrol to other zone?
	if plr.countUnitsInPatrolZone() > 0 {
		index := rnd.SelectRandomIndexFromWeighted(plr.countUnitsInPatrolZone(), func(x int) int {
			if plr.patrolZone[x] == nil {
				return 0
			}
			if plr.patrolZone[x].tapped {
				return 10
			}
			return plr.patrolZone[x].wounds
		})
		if index == -1 {
			// nothing.
		} else {
			unitToMove := plr.patrolZone[index]
			plr.moveUnit(unitToMove, PLAYERZONE_PATROL, index, PLAYERZONE_OTHER, 0)
			log.Printf("I moved %s (tapped: %v, wounds %d) to other zone.\n", unitToMove.getName(), unitToMove.tapped, unitToMove.wounds)
			moved = true
		}
	}
	// Maybe move something from other zone to patrol?
	if len(plr.otherZone) > 0 {
		index := rnd.SelectRandomIndexFromWeighted(len(plr.otherZone), func(x int) int {
			if plr.otherZone[x].tapped {
				return 0
			}
			if plr.otherZone[x].hasPassiveAbility(UPA_HEALING) {
				return 0
			}
			_, hp := plr.otherZone[x].getAtkHpWithWounds()
			return hp
		})
		if index == -1 {
			// nothing.
		} else {
			unitToMove := plr.otherZone[index]
			// select patrol place
			patrolIndex := rnd.SelectRandomIndexFromWeighted(5, func(x int) int {
				atk, hp := unitToMove.getAtkHpWithWounds()
				prob := 0
				switch x {
				case 0:
					prob = 5
					if hp <= 2 {
						prob = 6
					}
				case 1:
					prob = 2
					if atk <= 2 {
						prob = 4
					}
				case 2:
					prob = 2
					if plr.gold < 2 {
						prob = 4
					}
				case 3:
					prob = 2
					if plr.hand.size() < 3 {
						prob = 5 - plr.hand.size()
					}
				case 4:
					prob = 2
				}
				if plr.patrolZone[x] != nil {
					prob /= 2
				}
				return prob
			})
			plr.moveUnit(unitToMove, PLAYERZONE_OTHER, index, PLAYERZONE_PATROL, patrolIndex)
			log.Printf("I moved %s (tapped: %v, wounds %d) to patrol zone.\n", unitToMove.getName(), unitToMove.tapped, unitToMove.wounds)
			moved = true
		}
	}
	return moved
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

func (ai *aiPlayerController) tryAttack(g *game) bool {
	var candidates []*unit
	for _, u := range ai.controlsPlayer.otherZone {
		if g.canUnitAttack(u) {
			candidates = append(candidates, u)
		}
	}
	for _, u := range ai.controlsPlayer.patrolZone {
		if u != nil && g.canUnitAttack(u) {
			candidates = append(candidates, u)
		}
	}
	if len(candidates) == 0 {
		log.Println("I have nothing to attack with.")
		return false
	} else {
		attackerIndex := rnd.SelectRandomIndexFromWeighted(len(candidates), func(ind int) int {
			atk, hp := candidates[ind].getAtkHpWithWounds()
			if atk < 2 {
				return atk
			}
			return 3*atk + hp
		})
		attacker := candidates[attackerIndex]
		log.Printf("I attack with %s.", attacker.getName())
		return g.tryAttackAsUnit(ai.controlsPlayer, attacker)
	}
}

func (ai *aiPlayerController) actCodex(g *game) {
	plr := ai.controlsPlayer
	if plr.workers > 10 && !rnd.OneChanceFrom(3) {
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
