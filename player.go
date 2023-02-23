package main

type player struct {
	hand    cardStack
	draw    cardStack
	discard cardStack
	codices [3]codex // one per element

	patrolZone         [5]*unit
	commandZone        [3]*heroCard
	otherZone          []*unit
	cardsToAddNextTurn [2]card

	gold                int
	workers             int
	hiredWorkerThisTurn bool

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

func (p *player) sortHand() {
	p.hand.sortByName()
	p.hand.sortByCost()
}

func (p *player) hasTechLevel(lvl int) bool {
	if lvl == 0 {
		return true
	}
	for _, b := range p.techBuildings {
		if b.static.givesTech == lvl {
			return true
		}
	}
	return false
}

func (p *player) drawCard() {
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

const (
	PLAYERZONE_OTHER = iota
	PLAYERZONE_PATROL
)
