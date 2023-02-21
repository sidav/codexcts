package main

type player struct {
	hand    cardStack
	draw    cardStack
	discard cardStack
	// heroes  [3]*card need a separate type
	gold       int
	workers    int
	baseHealth int
}

func (p *player) sortHand() {
	p.hand.sortByCost()
}

func (p *player) drawCard() {
	poppedCard := p.draw.pop()
	p.hand.pushOnTop(poppedCard)
}
