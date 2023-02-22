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

func (g *game) canPlayerBuild(p *player, b *building) bool {
	if !b.static.isAddon {
		for _, tb := range p.techBuildings {
			if tb != nil && tb.static == b.static {
				return false
			}
		}
	} else if p.addonBuilding != nil {
		return false
	}
	return p.gold >= b.static.cost && p.workers >= b.static.requiresWorkers
}
