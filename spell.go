package main

type spell struct {
	isContinious     bool
	targetsNumber    int
	targetsAllInZone bool
	canTargetUnits   bool
	canTargetHeroes  bool
	canTargetOwner   bool
	canTargetEnemy   bool
	targetableZones  []int

	effectCode                 spellEffectCode
	effectValue1, effectValue2 int
}

type spellEffectCode int

const (
	SPELL_DEAL_DAMAGE = iota
	SPELL_ADD_RUNES
)
