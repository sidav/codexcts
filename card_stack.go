package main

import (
	"codexcts/lib/random"
	"fmt"
	"sort"
)

type cardStack []card

func (s *cardStack) pushOnTop(c card) {
	*s = append([]card{c}, *s...)
}

func (s *cardStack) addToBottom(c card) {
	*s = append(*s, c)
}

func (s *cardStack) pop() card {
	fmt.Printf("DEBUG: %d\n", len(*s))
	c := (*s)[0]
	*s = (*s)[1:]
	return c
}

func (s *cardStack) sortByCost() {
	sort.Slice(*s, func(i, j int) bool { return (*s)[i].getCost() < (*s)[j].getCost() })
}

func (s cardStack) shuffle(rnd random.PRNG) {
	// Fisher–Yates shuffle
	for i := len(s) - 1; i > 0; i-- {
		exchInd := rnd.Rand(i + 1)
		t := s[exchInd]
		s[exchInd] = s[i]
		s[i] = t
	}
}
