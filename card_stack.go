package main

import (
	"codexcts/lib/random"
	"sort"
)

type cardStack []card

func (s *cardStack) size() int {
	return len(*s)
}

func (s *cardStack) removeThis(c card) {
	for i := range *s {
		if (*s)[i] == c {
			*s = append((*s)[:i], (*s)[i+1:]...)
			return
		}
	}
	panic("No card " + c.getName() + "in stack!")
}

func (s *cardStack) pushOnTop(c card) {
	*s = append([]card{c}, *s...)
}

func (s *cardStack) addToBottom(c card) {
	*s = append(*s, c)
}

func (s *cardStack) pop() card {
	// fmt.Printf("DEBUG: %d\n", len(*s))
	c := (*s)[0]
	*s = (*s)[1:]
	return c
}

func (s *cardStack) moveFrom(s2 *cardStack) {
	s.addToBottom(s2.pop())
}

func (s *cardStack) sortByCost() {
	sort.Slice(*s, func(i, j int) bool { return (*s)[i].getCost() < (*s)[j].getCost() })
}

func (s *cardStack) sortByName() {
	sort.Slice(*s, func(i, j int) bool { return (*s)[i].getName()[0] < (*s)[j].getName()[0] })
}

func (s *cardStack) sortByType() {
	sort.Slice(*s, func(i, j int) bool {
		switch (*s)[i].(type) {
		case *unitCard:
			return true
		case *magicCard:
			switch (*s)[j].(type) {
			case *unitCard:
				return false
			case *magicCard:
				return true
			}
		}
		return false
	})
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
