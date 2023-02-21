package main

type player struct {
	hand    cardStack
	draw    cardStack
	discard cardStack

	patrolZone         [5]*unit
	commandZone        [3]*heroCard
	otherZone          []*unit
	cardsToAddNextTurn [2]card

	gold    int
	workers int

	baseHealth    int
	techBuildings [3]*building
	addonBuilding *building
}

func (p *player) sortHand() {
	p.hand.sortByName()
	p.hand.sortByCost()
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
