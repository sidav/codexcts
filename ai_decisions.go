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
	action := rnd.SelectRandomIndexFromWeighted(AiDecisionsCount, func(act int) int {
		switch act {
		case AiDecisionPlayUnit:
			return 2
		case AiDecisionPlayHero:
			return 2
		case AiDecisionPlayMagic:
			if ai.controlsPlayer.hasHeroOnField() {
				return 1
			}
			return 0
		case AiDecisionBuild:
			return 1
		case AiDecisionLevelUp:
			return 2
		case AiDecisionAbstain:
			return 1
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
