package main

type player struct {
	hand    cardStack
	draw    cardStack
	discard cardStack

	patrolZone  [5]*creature
	commandZone [3]*heroCard
	gold        int
	workers     int
	baseHealth  int
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
