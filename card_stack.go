package main

import "codexcts/lib/random"

type cardStack []card

func (s *cardStack) pushOnTop(c card) {
	*s = append([]card{c}, *s...)
}

func (s *cardStack) addToBottom(c card) {
	*s = append(*s, c)
}

func (s *cardStack) pop() card {
	c := (*s)[0]
	*s = (*s)[1:]
	return c
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
