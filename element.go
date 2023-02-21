package main

type element uint8

const (
	ELEMENT_NEUTRAL element = iota
	ELEMENT_BASHING
	ELEMENT_FINESSE
)

func (e element) getName() string {
	switch e {
	case ELEMENT_NEUTRAL:
		return "Neutral"
	case ELEMENT_FINESSE:
		return "Finesse"
	case ELEMENT_BASHING:
		return "Bashing"
	default:
		panic("No such element for getName")
	}
}
