package main

import "log"

func (ai *aiPlayerController) tryMoveUnits(g *game) bool {
	plr := ai.controlsPlayer
	moved := false
	// Maybe move something from patrol to other zone?
	index := rnd.SelectRandomIndexFromWeighted(5, func(x int) int {
		if plr.patrolZone[x] == nil {
			return 0
		}
		if plr.patrolZone[x].tapped {
			return 10
		}
		return plr.patrolZone[x].wounds
	})
	if index == -1 {
		log.Println("I decided not to remove any patrollers.")
	} else {
		unitToMove := plr.patrolZone[index]
		plr.moveUnit(unitToMove, PLAYERZONE_PATROL, index, PLAYERZONE_OTHER, 0)
		log.Printf("I moved %s (tapped: %v, wounds %d) to other zone.\n", unitToMove.getName(), unitToMove.tapped, unitToMove.wounds)
		moved = true
	}
	// Maybe move something from other zone to patrol?
	if plr.countUnitsInPatrolZone() >= 2 && rnd.Rand(4) > 0 {
		log.Println("I think my patrol zone is full enough.")
		return moved
	}
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

func (ai *aiPlayerController) tryAttack(g *game) bool {
	var candidates []*unit
	usePatrollers := rnd.OneChanceFrom(2) ||
		g.getEnemyForPlayer(ai.controlsPlayer).countUntappedUnitsInPatrolZone() <= ai.controlsPlayer.countUntappedUnitsInPatrolZone()
	for _, u := range ai.controlsPlayer.otherZone {
		if g.canUnitAttack(u) {
			candidates = append(candidates, u)
		}
	}
	if usePatrollers {
		for _, u := range ai.controlsPlayer.patrolZone {
			if u != nil && g.canUnitAttack(u) {
				candidates = append(candidates, u)
			}
		}
	}
	if len(candidates) == 0 {
		log.Println("I have nothing to attack with.")
		return false
	} else {
		attackerIndex := rnd.SelectRandomIndexFromWeighted(len(candidates), func(ind int) int {
			atk, hp := candidates[ind].getAtkHpWithWounds()
			if candidates[ind].hasPassiveAbility(UPA_HEALING) {
				return 0
			}
			if candidates[ind].hasPassiveAbility(UPA_HASTE) {
				return 2 * atk
			}
			if candidates[ind].hasPassiveAbility(UPA_FRENZY) {
				return 3*(atk+1) + hp
			}
			if atk < 2 {
				return atk
			}
			return 3*atk + hp
		})
		if attackerIndex == -1 {
			log.Println("I decided not to attack.")
			return false
		}
		attacker := candidates[attackerIndex]
		log.Printf("I attack with %s.", attacker.getName())
		return g.tryAttackAsUnit(ai.controlsPlayer, attacker)
	}
}
