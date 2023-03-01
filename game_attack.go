package main

import "fmt"

func (g *game) getAttackableCoordsForUnit(attacker *unit, owner *player) []*playerZoneCoords {
	enemy := g.getEnemyForPlayer(owner)
	attackerFlying := attacker.hasPassiveAbility(UPA_FLYING)
	attackerAA := attacker.hasPassiveAbility(UPA_ANTI_AIR)
	if enemy.patrolZone[0] != nil {
		if attackerFlying == enemy.patrolZone[0].hasPassiveAbility(UPA_FLYING) {
			return []*playerZoneCoords{{enemy, PLAYERZONE_PATROL, 0}}
		}
	}
	onlyPatrolZone := false
	list := make([]*playerZoneCoords, 0)
	for i, patroller := range enemy.patrolZone {
		if patroller != nil {
			if attackerAA == patroller.hasPassiveAbility(UPA_FLYING) || attackerFlying == patroller.hasPassiveAbility(UPA_FLYING) {
				list = append(list, &playerZoneCoords{enemy, PLAYERZONE_PATROL, i})
				if !patroller.tapped && attackerFlying == patroller.hasPassiveAbility(UPA_FLYING) {
					onlyPatrolZone = true
				}
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
				if attackerAA == p.hasPassiveAbility(UPA_FLYING) || attackerFlying == p.hasPassiveAbility(UPA_FLYING) {
					list = append(list, &playerZoneCoords{enemy, PLAYERZONE_OTHER, i})
				}
			}
		}
	}
	return list
}

func (g *game) resolveAttack(attacker *unit, attackingPlayer *player, targetCoords *playerZoneCoords) {
	atk, _ := attacker.getAtkHpWithWounds()
	defendingPlayer := targetCoords.player
	var targetUnit *unit
	targetIndex := targetCoords.indexInZone
	targetArmorBonus := 0
	targetAttackBonus := 0

	if defendingPlayer.addonBuilding != nil && defendingPlayer.addonBuilding.static.damagesAttackers {
		g.messageForPlayer += fmt.Sprintf("%s's %s took 1 damage from %s's tower. \n ", attackingPlayer.name,
			attacker.getName(), defendingPlayer.name)
		attacker.wounds++
	}

	switch targetCoords.zone {
	case PLAYERZONE_MAIN_BASE:
		defendingPlayer.baseHealth -= atk
		g.messageForPlayer += fmt.Sprintf("%s's base took %d damage. (%d HP remaining)\n ", defendingPlayer.name,
			atk, defendingPlayer.baseHealth)
	case PLAYERZONE_TECH_BUILDINGS:
		defendingPlayer.techBuildings[targetIndex].currentHitpoints -= atk
		g.messageForPlayer += fmt.Sprintf("%s's %s took %d damage (%d HP remaining). \n ", defendingPlayer.name,
			defendingPlayer.techBuildings[targetIndex].static.name, atk, defendingPlayer.techBuildings[targetIndex].currentHitpoints)
		if defendingPlayer.techBuildings[targetIndex].currentHitpoints <= 0 {
			defendingPlayer.baseHealth -= 2
			g.messageForPlayer += fmt.Sprintf(" It's destroyed! %s's base took 2 damage (%d HP remaining). \n ",
				defendingPlayer.name, defendingPlayer.baseHealth)
			defendingPlayer.techBuildings[targetIndex] = nil // TODO: rebuildability
		}
	case PLAYERZONE_ADDON_BUILDING:
		defendingPlayer.addonBuilding.currentHitpoints -= atk
		g.messageForPlayer += fmt.Sprintf("%s's %s took %d damage (%d HP remaining). \n ", defendingPlayer.name,
			defendingPlayer.addonBuilding.static.name, atk, defendingPlayer.addonBuilding.currentHitpoints)
		if defendingPlayer.addonBuilding.currentHitpoints <= 0 {
			defendingPlayer.baseHealth -= 2
			g.messageForPlayer += fmt.Sprintf(" It's destroyed! %s's base took 2 damage (%d HP remaining). \n ",
				defendingPlayer.name, defendingPlayer.baseHealth)
			defendingPlayer.addonBuilding = nil
		}
	case PLAYERZONE_OTHER:
		targetUnit = defendingPlayer.otherZone[targetIndex]
	case PLAYERZONE_PATROL:
		targetUnit = defendingPlayer.patrolZone[targetIndex]
		switch targetIndex {
		case 0:
			if defendingPlayer.patrolLeaderHasShield {
				g.messageForPlayer += fmt.Sprintf("%s's %s is in leader slot, thus getting 1 armor. \n ",
					defendingPlayer.name,
					targetUnit.getName())
				targetArmorBonus++
				defendingPlayer.patrolLeaderHasShield = false
			} else {
				g.messageForPlayer += "Defender has already spent its shield from patrol position. \n "
			}
		case 1:
			g.messageForPlayer += fmt.Sprintf("%s's %s is in elite slot, thus getting +1 damage. \n ",
				defendingPlayer.name,
				targetUnit.getName())
			targetAttackBonus++
		}
	}
	// dealing the damage to unit
	if targetUnit != nil {
		g.messageForPlayer += fmt.Sprintf("Attacker: %s \n ", attacker.getNameWithStats())
		if attacker.hasPassiveAbility(UPA_FRENZY) {
			inc := attacker.getPassiveAbilityValue(UPA_FRENZY)
			g.messageForPlayer += fmt.Sprintf("  Attacker has Frenzy %d (+%d ATK) \n ", inc, inc)
			atk += inc
		}
		g.messageForPlayer += fmt.Sprintf("Defender: %s \n ", targetUnit.getNameWithStats())
		atk -= targetArmorBonus
		backAtk, _ := targetUnit.getAtkHpWithWounds()
		backAtk += targetAttackBonus
		targetUnit.wounds += atk
		g.messageForPlayer += fmt.Sprintf("Defending %s took %d damage from attacker, now having %d wounds. \n ",
			targetUnit.getName(), atk, targetUnit.wounds)
		if attacker.hasPassiveAbility(UPA_FLYING) && !(targetUnit.hasPassiveAbility(UPA_FLYING) || targetUnit.hasPassiveAbility(UPA_ANTI_AIR)) {
			g.messageForPlayer += fmt.Sprintf("Attacking %s took no damage from defender, as it has flying. \n ", attacker.getName())
		} else {
			attacker.wounds += backAtk
			g.messageForPlayer += fmt.Sprintf("Attacking %s took %d damage from defender, now having %d wounds. \n ",
				attacker.getName(), backAtk, attacker.wounds)
		}
	}
}

func (g *game) removeDeadUnits() {
	for _, p := range g.players {
		for ind := len(p.otherZone) - 1; ind >= 0; ind-- {
			unt := p.otherZone[ind]
			_, hp := unt.getAtkHpWithWounds()
			if hp <= 0 {
				g.messageForPlayer += fmt.Sprintf("%s's %s dies. \n ", p.name, unt.getName())
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
			_, hp := unt.getAtkHpWithWounds()
			if hp <= 0 {
				g.messageForPlayer += fmt.Sprintf("%s's %s dies. \n ", p.name, unt.getName())
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

func (g *game) dealDamageByCoords(damageAmount int, targetCoords *playerZoneCoords) {
	targetPlayer := targetCoords.player
	targetIndex := targetCoords.indexInZone
	targetArmor := 0
	var targetUnit *unit
	switch targetCoords.zone {
	case PLAYERZONE_MAIN_BASE:
		targetCoords.player.baseHealth -= damageAmount
		g.messageForPlayer += fmt.Sprintf("%s's base took %d damage. (%d HP remaining)\n ", targetCoords.player.name,
			damageAmount, targetCoords.player.baseHealth)
	case PLAYERZONE_TECH_BUILDINGS:
		targetPlayer.techBuildings[targetIndex].currentHitpoints -= damageAmount
		g.messageForPlayer += fmt.Sprintf("%s's %s took %d damage (%d HP remaining). \n ", targetPlayer.name,
			targetPlayer.techBuildings[targetIndex].static.name, damageAmount, targetPlayer.techBuildings[targetIndex].currentHitpoints)
		if targetPlayer.techBuildings[targetIndex].currentHitpoints <= 0 {
			targetPlayer.baseHealth -= 2
			g.messageForPlayer += fmt.Sprintf(" It's destroyed! %s's base took 2 damage (%d HP remaining). \n ",
				targetPlayer.name, targetPlayer.baseHealth)
			targetPlayer.techBuildings[targetIndex] = nil // TODO: rebuildability
		}
	case PLAYERZONE_ADDON_BUILDING:
		targetPlayer.addonBuilding.currentHitpoints -= damageAmount
		g.messageForPlayer += fmt.Sprintf("%s's %s took %d damage (%d HP remaining). \n ", targetPlayer.name,
			targetPlayer.addonBuilding.static.name, damageAmount, targetPlayer.addonBuilding.currentHitpoints)
		if targetPlayer.addonBuilding.currentHitpoints <= 0 {
			targetPlayer.baseHealth -= 2
			g.messageForPlayer += fmt.Sprintf(" It's destroyed! %s's base took 2 damage (%d HP remaining). \n ",
				targetPlayer.name, targetPlayer.baseHealth)
			targetPlayer.addonBuilding = nil
		}
	case PLAYERZONE_OTHER:
		targetUnit = targetPlayer.otherZone[targetIndex]
	case PLAYERZONE_PATROL:
		targetUnit = targetPlayer.patrolZone[targetIndex]
		if targetIndex == 0 {
			if targetPlayer.patrolLeaderHasShield {
				g.messageForPlayer += fmt.Sprintf("%s's %s is in leader slot, thus getting 1 armor. \n ",
					targetPlayer.name,
					targetUnit.getName())
				targetArmor++
				targetPlayer.patrolLeaderHasShield = false
			} else {
				g.messageForPlayer += "Defender has already spent its shield from patrol position. \n "
			}
		}
	}
	if targetUnit != nil {
		finalDmg := damageAmount - targetArmor
		g.messageForPlayer += fmt.Sprintf("%s receives %d damage.", targetUnit.getName(), finalDmg)
		targetUnit.wounds += finalDmg
	}
}
