package main

func (g *game) tryPlayCardAsWorker(c card) bool {
	if g.currentPlayer.gold > 0 && !g.currentPlayer.hiredWorkerThisTurn {
		g.currentPlayer.hand.removeThis(c)
		g.currentPlayer.workers++
		g.currentPlayer.gold--
		g.currentPlayer.hiredWorkerThisTurn = true
		return true
	}
	return false
}

func (g *game) tryPlayUnitCardFromHand(c card) bool {
	if g.canPlayerPlayCard(g.currentPlayer, c) {
		unt := &unit{
			card:   c,
			tapped: false,
			wounds: 0,
		}
		unt.attackedThisTurn = !unt.hasPassiveAbility(UPA_HASTE)
		g.currentPlayer.otherZone = append(g.currentPlayer.otherZone, unt)
		g.currentPlayer.hand.removeThis(c)
		g.currentPlayer.gold -= c.getCost()
		return true
	}
	return false
}

func (g *game) tryPlayHeroCard(c card) bool {
	if g.canPlayerPlayCard(g.currentPlayer, c) {
		unt := &unit{
			card:             c,
			tapped:           false,
			attackedThisTurn: true, // so that it can't attack just now
			wounds:           0,
			level:            1,
		}
		unt.attackedThisTurn = !unt.hasPassiveAbility(UPA_HASTE)
		g.currentPlayer.otherZone = append(g.currentPlayer.otherZone, unt)
		g.currentPlayer.gold -= c.getCost()
		for i, h := range g.currentPlayer.commandZone {
			if h == c {
				g.currentPlayer.commandZone[i] = nil
				return true
			}
		}
		panic("Something is wrong when playing a hero")
	}
	return false
}

func (g *game) tryBuildNextTechForPlayer(p *player) bool {
	for i := range p.techBuildings {
		if p.techBuildings[i] == nil {
			tb := getTechBuildingByTechLevel(i + 1)
			g.tryBuildBuildingForPlayer(p, tb)
			return true
		}
	}
	return false
}

func (g *game) tryLevelUpHero(p *player, unit *unit) bool {
	if unit.isHero() {
		if p.gold > 0 && unit.card.(*heroCard).getMaxLevel() > unit.level {
			p.gold--
			unit.level++
			// heal hero on a threshold
			for _, lad := range unit.card.(*heroCard).levelsAttDef {
				if lad[0] == unit.level {
					unit.wounds = 0
				}
			}
			return true
		}
	}
	return false
}

func (g *game) tryBuildBuildingForPlayer(p *player, b *buildingStatic) bool {
	if !g.canPlayerBuild(p, b) {
		return false
	}
	if b.isAddon {
		p.addonBuilding = &building{
			static:              b,
			currentHitpoints:    b.maxHitpoints,
			isUnderConstruction: true,
		}
	} else {
		p.techBuildings[b.givesTech-1] = &building{
			static:              b,
			currentHitpoints:    b.maxHitpoints,
			isUnderConstruction: true,
		}
	}
	return true
}

func (g *game) tryAddCardFromCodex(p *player, c card, codexIndex int) bool {
	if p.codices[codexIndex].getCardCount(c) == 0 {
		return false
	}
	p.codices[codexIndex].removeSingleCard(c)
	p.cardsToAddNextTurn[p.cardsAddedFromCodexThisTurn] = c
	p.cardsAddedFromCodexThisTurn++
	return true
}
