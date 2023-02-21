package main

import (
	"codexcts/lib/random"
	"codexcts/lib/random/pcgrandom"
	"fmt"
)

func main() {
	var rnd random.PRNG = pcgrandom.New(-1)
	var deck cardStack
	sc := getStartingCardsForElement(ELEMENT_NEUTRAL)
	fmt.Printf("Got %d cards in codex\n", len(sc))
	for _, c := range sc {
		deck.addToBottom(c)
	}
	fmt.Printf("Got %d cards in deck:\n", len(deck))
	for _, c := range deck {
		fmt.Printf(" %s\n", c.getName())
	}
	deck.shuffle(rnd)
	fmt.Printf("Shuffled:\n")
	for _, c := range deck {
		fmt.Printf(" %s\n", c.getName())
	}
	fmt.Printf("Hand:\n")
	var hand cardStack
	for i := 0; i < 5; i++ {
		hand.pushOnTop(deck.pop())
	}
	for _, c := range hand {
		fmt.Printf(" %s\n", c.getName())
	}
	fmt.Printf("Deck after draw:\n")
	for _, c := range deck {
		fmt.Printf(" %s\n", c.getName())
	}
}
