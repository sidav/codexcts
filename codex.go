package main

type codex struct {
	cards       cardStack
	cardsCounts []int
}

func (cd *codex) getRemainingUniqueCardsCount() int {
	count := 0
	for i := range cd.cards {
		if cd.cardsCounts[i] > 0 {
			count++
		}
	}
	return count
}

func (cd *codex) getTotalCardsCount() int {
	count := 0
	for i := range cd.cards {
		count += cd.cardsCounts[i]
	}
	return count
}

func (cd *codex) getCardByIndex(index int) card {
	return cd.cards[index]
}

func (cd *codex) getCardCount(c card) int {
	for i := range cd.cards {
		if cd.cards[i] == c {
			return cd.cardsCounts[i]
		}
	}
	panic("No such card in codex!")
}

func (cd *codex) removeSingleCard(c card) {
	for i := range cd.cards {
		if cd.cards[i] == c && cd.cardsCounts[i] > 0 {
			cd.cardsCounts[i]--
			return
		}
	}
	panic("No such card in codex!")
}

func (cd *codex) addCard(c card) {
	for i := range cd.cards {
		if cd.cards[i] == c {
			cd.cardsCounts[i]++
			return
		}
	}
	cd.cards = append(cd.cards, c)
	cd.cardsCounts = append(cd.cardsCounts, 1)
}
