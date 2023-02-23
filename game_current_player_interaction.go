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
	if g.currentPlayer.gold >= c.getCost() {
		g.currentPlayer.otherZone = append(g.currentPlayer.otherZone, &unit{
			card:   c,
			tapped: false,
			wounds: 0,
			level:  0,
		})
		g.currentPlayer.hand.removeThis(c)
		g.currentPlayer.gold -= c.getCost()
		return true
	}
	return false
}

func (g *game) canPlayerBuild(p *player, b *buildingStatic) bool {
	if !b.isAddon {
		for _, tb := range p.techBuildings {
			if tb != nil && tb.static == b {
				return false
			}
		}
	} else if p.addonBuilding != nil {
		return false
	}
	return p.gold >= b.cost && p.workers >= b.requiresWorkers
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
		p.techBuildings[b.givesTech] = &building{
			static:              b,
			currentHitpoints:    b.maxHitpoints,
			isUnderConstruction: true,
		}
	}
	return true
}
