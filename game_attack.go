package main

func (g *game) getAttackableCoordsForUnit(attacker *unit, owner *player) []*playerZoneCoords {
	enemy := g.getEnemyForPlayer(owner)
	if enemy.patrolZone[0] != nil {
		return []*playerZoneCoords{{enemy, PLAYERZONE_PATROL, 0}}
	}
	onlyPatrolZone := false
	list := make([]*playerZoneCoords, 0)
	for i, p := range enemy.patrolZone {
		if p != nil {
			list = append(list, &playerZoneCoords{enemy, PLAYERZONE_PATROL, i})
			if !p.tapped {
				onlyPatrolZone = true
			}
		}
	}
	if !onlyPatrolZone { // patrol zone empty, adding everything the player has
		list = append(list, &playerZoneCoords{enemy, PLAYERZONE_MAIN_BASE, 0})
		for i, t := range enemy.techBuildings {
			if t != nil {
				list = append(list, &playerZoneCoords{enemy, PLAYERZONE_TECH_BUILDINGS, i})
			}
		}
		if enemy.addonBuilding != nil {
			list = append(list, &playerZoneCoords{enemy, PLAYERZONE_ADDON_BUILDING, 0})
		}
		for i, p := range enemy.otherZone {
			if p != nil {
				list = append(list, &playerZoneCoords{enemy, PLAYERZONE_OTHER, i})
			}
		}
	}
	return list
}

func (g *game) tryAttackAsUnit(owner *player, attacker *unit) bool {
	if !g.canUnitAttack(attacker) {
		return false
	}
	coords := g.getAttackableCoordsForUnit(attacker, owner)
	if len(coords) == 0 {
		return false
	}
	var selectedCoords *playerZoneCoords
	if len(coords) == 1 {
		selectedCoords = coords[0]
	} else {
		selectedCoords = g.playersControllers[g.getPlayerNumber(owner)].selectCoordsFromListCallback(
			"Select the target of the attack", coords)
	}
	g.performAttack(attacker, owner, selectedCoords)
	g.removeDeadUnits()
	attacker.tapped = true
	return true
}

func (g *game) performAttack(attacker *unit, attackingPlayer *player, targetCoords *playerZoneCoords) {
	atk, _ := attacker.getAtkHp()
	defendingPlayer := targetCoords.player
	var targetUnit *unit
	targetIndex := targetCoords.indexInZone
	targetArmorBonus := 0
	targetAttackBonus := 0

	switch targetCoords.zone {
	case PLAYERZONE_MAIN_BASE:
		defendingPlayer.baseHealth -= atk
	case PLAYERZONE_TECH_BUILDINGS:
		defendingPlayer.techBuildings[targetIndex].currentHitpoints -= atk
		if defendingPlayer.techBuildings[targetIndex].currentHitpoints <= 0 {
			defendingPlayer.baseHealth -= 2
		}
	case PLAYERZONE_ADDON_BUILDING:
		defendingPlayer.addonBuilding.currentHitpoints -= atk
		if defendingPlayer.addonBuilding.currentHitpoints <= 0 {
			defendingPlayer.baseHealth -= 2
		}
	case PLAYERZONE_OTHER:
		targetUnit = defendingPlayer.otherZone[targetIndex]
	case PLAYERZONE_PATROL:
		targetUnit = defendingPlayer.patrolZone[targetIndex]
		switch targetIndex {
		case 0:
			targetArmorBonus++
		case 1:
			targetAttackBonus++
		}
	}
	// dealing the damage to unit
	if targetUnit != nil {
		backAtk, _ := targetUnit.getAtkHp()
		targetUnit.wounds += atk - targetArmorBonus
		attacker.wounds += backAtk + targetAttackBonus
	}
	if defendingPlayer.addonBuilding != nil && defendingPlayer.addonBuilding.static.damagesAttackers {
		attacker.wounds++
	}
}

func (g *game) removeDeadUnits() {
	for _, p := range g.players {
		for ind := len(p.otherZone) - 1; ind >= 0; ind-- {
			unt := p.otherZone[ind]
			_, hp := unt.getAtkHp()
			if unt.wounds >= hp {
				if unt.isHero() {
					for heroInd := range p.commandZone {
						if p.commandZone[heroInd] == nil {
							p.commandZone[heroInd] = unt.card.(*heroCard)
							break
						}
					}
				} else {
					p.discard.addToBottom(unt.card)
				}
				p.otherZone = append(p.otherZone[:ind], p.otherZone[ind+1:]...)
			}
		}
		for ind := range p.patrolZone {
			unt := p.patrolZone[ind]
			if unt == nil {
				continue
			}
			_, hp := unt.getAtkHp()
			if unt.wounds >= hp {
				if unt.isHero() {
					for heroInd := range p.commandZone {
						if p.commandZone[heroInd] == nil {
							p.commandZone[heroInd] = unt.card.(*heroCard)
							break
						}
					}
				} else {
					p.discard.addToBottom(unt.card)
				}
				p.patrolZone[ind] = nil
				switch ind {
				case 2:
					p.gold++
				case 3:
					p.drawCard()
				}
			}
		}
	}
}
