package main

func (g *game) getAttackableCoordsForUnit(attacker *unit, owner *player) []*playerZoneCoords {
	enemy := g.getEnemyForPlayer(owner)
	if enemy.patrolZone[0] != nil {
		return []*playerZoneCoords{{enemy, PLAYERZONE_PATROL, 0}}
	}
	list := make([]*playerZoneCoords, 0)
	for i, p := range enemy.patrolZone {
		if p != nil {
			list = append(list, &playerZoneCoords{enemy, PLAYERZONE_PATROL, i})
		}
	}
	if len(list) == 0 { // patrol zone empty, adding everything the player has
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
	if attacker.tapped {
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
	return true
}

func (g *game) performAttack(attacker *unit, attackerOwner *player, targetCoords *playerZoneCoords) {
	atk, _ := attacker.getAtkHp()
	targetOwner := targetCoords.player
	switch targetCoords.zone {
	case PLAYERZONE_MAIN_BASE:
		targetOwner.baseHealth -= atk
	case PLAYERZONE_TECH_BUILDINGS:
		targetOwner.techBuildings[targetCoords.indexInZone].currentHitpoints -= atk
		if targetOwner.techBuildings[targetCoords.indexInZone].currentHitpoints <= 0 {
			targetOwner.baseHealth -= 2
		}
	case PLAYERZONE_ADDON_BUILDING:
		targetOwner.addonBuilding.currentHitpoints -= atk
		if targetOwner.addonBuilding.currentHitpoints <= 0 {
			targetOwner.baseHealth -= 2
		}
	case PLAYERZONE_OTHER:
		target := targetOwner.otherZone[targetCoords.indexInZone]
		target.wounds += atk
	case PLAYERZONE_PATROL:
		target := targetOwner.patrolZone[targetCoords.indexInZone]
		backAtk, _ := target.getAtkHp()
		target.wounds += atk
		attacker.wounds += backAtk
	}
}
