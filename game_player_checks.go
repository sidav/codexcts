package main

func (g *game) canUnitAttack(u *unit) bool {
	return !(u.tapped || u.attackedThisTurn)
}

func (g *game) canPlayerPlayCard(p *player, c card) bool {
	var can bool
	switch c.(type) {
	case *magicCard:
		can = false // TODO: check for hero presence etc
	case *unitCard:
		can = p.hasTechLevel(c.(*unitCard).techLevel)
	case *heroCard:
		can = true
		// TODO: check for heroes number etc
	}
	return can && p.gold >= c.getCost()
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
