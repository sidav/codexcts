package main

import "log"

const (
	AiDecisionPlayUnit = iota
	AiDecisionPlayHero
	AiDecisionPlayMagic
	AiDecisionBuild
	AiDecisionLevelUp
	AiDecisionAbstain
	AiDecisionsCount
)

func (ai *aiPlayerController) decideRandomEconomicAction(g *game) int {
	if ai.controlsPlayer.gold == 0 {
		return AiDecisionAbstain
	}

	const baseWeight = 3
	weightForSpendingCard := ai.controlsPlayer.hand.size()
	
	action := rnd.SelectRandomIndexFromWeighted(AiDecisionsCount, func(act int) int {
		switch act {
		case AiDecisionPlayUnit:
			return weightForSpendingCard
		case AiDecisionPlayHero:
			return baseWeight
		case AiDecisionPlayMagic:
			if ai.controlsPlayer.hasHeroOnField() {
				return weightForSpendingCard
			}
			return 0
		case AiDecisionBuild:
			return baseWeight
		case AiDecisionLevelUp:
			return baseWeight
		case AiDecisionAbstain:
			return baseWeight
		default:
			panic("Wat")
		}
	})
	return action
}

func (ai *aiPlayerController) tryPerformRandomEconomicAction(g *game) bool {
	// TODO: build
	actionToPerform := ai.decideRandomEconomicAction(g)
	switch actionToPerform {
	case AiDecisionPlayUnit:
		return ai.tryPlayUnit(g)
	case AiDecisionPlayHero:
		return ai.playHero(g)
	case AiDecisionPlayMagic:
		return ai.tryPlayMagicCard(g)
	case AiDecisionBuild:
		return ai.tryBuild(g)
	case AiDecisionLevelUp:
		return ai.tryLevelUpHero(g)
	case AiDecisionAbstain:
		log.Println("I skipped an action.")
		return true
	}
	panic("Strange random in action selection")
}
