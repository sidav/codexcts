package main

import "log"

const (
	AiDecisionPlayUnit = iota
	AiDecisionPlayHero
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
			return 1
		case AiDecisionPlayHero:
			return 2
		case AiDecisionBuild:
			return 1
		case AiDecisionLevelUp:
			return 1
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
