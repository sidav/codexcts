package main

type player struct {
	name string

	hand    cardStack
	draw    cardStack
	discard cardStack
	codices [3]codex // one per element

	patrolZone                  [5]*unit
	commandZone                 [3]*heroCard
	otherZone                   []*unit
	cardsToAddNextTurn          [2]card
	cardsAddedFromCodexThisTurn int

	gold                  int
	workers               int
	hiredWorkerThisTurn   bool
	patrolLeaderHasShield bool

	baseHealth    int
	techBuildings [3]*building
	addonBuilding *building
}

func (p *player) moveUnit(u *unit, fromZone, indexFrom, toZone, indexTo int) {
	var swapWith *unit
	if toZone == PLAYERZONE_PATROL && p.patrolZone[indexTo] != nil {
		swapWith = p.patrolZone[indexTo]
	}
	if fromZone == PLAYERZONE_OTHER {
		if p.otherZone[indexFrom] != u {
			panic("Wat")
		}
		if toZone == PLAYERZONE_OTHER {
			return
		}
		p.otherZone = append(p.otherZone[:indexFrom], p.otherZone[indexFrom+1:]...)
		p.patrolZone[indexTo] = u
		if swapWith != nil {
			p.otherZone = append(p.otherZone, swapWith)
		}
	}
	if fromZone == PLAYERZONE_PATROL {
		p.patrolZone[indexFrom] = nil
		if toZone == PLAYERZONE_OTHER {
			p.otherZone = append(p.otherZone, u)
		}
		if toZone == PLAYERZONE_PATROL {
			p.patrolZone[indexTo] = u
			p.patrolZone[indexFrom] = swapWith
		}
	}
}

func (p *player) getUnitsInAllActiveZones() (units []*unit) {
	for _, u := range p.otherZone {
		units = append(units, u)
	}
	for _, u := range p.patrolZone {
		if u != nil {
			units = append(units, u)
		}
	}
	return
}

func (p *player) hasHeroOnField() bool {
	for _, u := range p.otherZone {
		if u.isHero() {
			return true
		}
	}
	for _, u := range p.patrolZone {
		if u != nil && u.isHero() {
			return true
		}
	}
	return false
}

func (p *player) countUnitsInPatrolZone() int {
	sum := 0
	for _, u := range p.patrolZone {
		if u != nil {
			sum++
		}
	}
	return sum
}

func (p *player) getUnitByCoords(coords *playerZoneCoords) *unit {
	switch coords.zone {
	case PLAYERZONE_OTHER:
		return p.otherZone[coords.indexInZone]
	case PLAYERZONE_PATROL:
		return p.patrolZone[coords.indexInZone]
	}
	panic("Strange coords given!")
}

func (p *player) countUntappedUnitsInPatrolZone() int {
	sum := 0
	for _, u := range p.patrolZone {
		if u != nil && !u.tapped {
			sum++
		}
	}
	return sum
}

func (p *player) sortHand() {
	p.hand.sortByName()
	p.hand.sortByCost()
}

func (p *player) hasTechLevel(lvl int) bool {
	if lvl == 0 {
		return true
	}
	for _, b := range p.techBuildings {
		if b != nil && !b.isUnderConstruction && b.static.givesTech == lvl {
			return true
		}
	}
	return false
}

func (p *player) drawCard() {
	if p.draw.size() == 0 {
		p.addDiscardIntoDraw()
		p.shuffleDraw()
	}
	if p.draw.size() == 0 {
		return
	}
	p.hand.moveFrom(&p.draw)
}

func (p *player) discardHand() {
	for len(p.hand) > 0 {
		p.discard.moveFrom(&p.hand)
	}
}

func (p *player) shuffleDraw() {
	p.draw.shuffle(rnd)
}

func (p *player) addDiscardIntoDraw() {
	for len(p.discard) > 0 {
		p.draw.moveFrom(&p.discard)
	}
}

func (p *player) isObligatedToAdd2Cards() bool {
	return p.workers < 10
}
