package main

import "log"

const (
	AiDecisionPlayUnit = iota
	AiDecisionPlayHero
	AiDecisionMoveUnit
	AiDecisionBuild
	AiDecisionLevelUp
	AiDecisionAttack
	AiDecisionAbstain
	AiDecisionsCount
)

func (ai *aiPlayerController) decideRandomAction(g *game) int {
	action := rnd.SelectRandomIndexFromWeighted(AiDecisionsCount, func(act int) int {
		switch act {
		case AiDecisionPlayUnit:
			return 1
		case AiDecisionPlayHero:
			return 2
		case AiDecisionMoveUnit:
			return len(ai.controlsPlayer.otherZone)
		case AiDecisionBuild:
			return 1
		case AiDecisionLevelUp:
			return 1
		case AiDecisionAttack:
			if g.getEnemyForPlayer(ai.controlsPlayer).countUnitsInPatrolZone() == 0 {
				return 20
			}
			return len(ai.controlsPlayer.getUnitsInAllActiveZones())
		case AiDecisionAbstain:
			return 1
		default:
			panic("Wat")
		}
	})
	return action
}

func (ai *aiPlayerController) tryPerformRandomAction(g *game) bool {
	// TODO: build
	actionToPerform := ai.decideRandomAction(g)
	switch actionToPerform {
	case AiDecisionPlayUnit:
		return ai.tryPlayUnit(g)
	case AiDecisionPlayHero:
		return ai.playHero(g)
	case AiDecisionMoveUnit:
		return ai.tryMoveUnit(g)
	case AiDecisionBuild:
		return ai.tryBuild(g)
	case AiDecisionLevelUp:
		return ai.tryLevelUpHero(g)
	case AiDecisionAttack:
		return ai.tryAttack(g)
	case AiDecisionAbstain:
		log.Println("I skipped an action.")
		return true
	}
	panic("Strange random in action selection")
}
