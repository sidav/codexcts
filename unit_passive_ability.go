package main

import (
	"fmt"
	"strconv"
)

type unitPassiveAbility struct {
	code               unitPassiveAbilityCode
	availableFromLevel int // heroes only
	value              int
}

func (us *unitPassiveAbility) getFormattedName() string {
	name := ""
	switch us.code {
	case UPA_HASTE:
		name = "Haste"
	//case UPA_HEALING:
	//	name = "Healing"
	case UPA_READINESS:
		name = "Readiness"
	//case UPA_FLYING:
	//	name = "Flying"
	case UPA_FRENZY:
		name = "Frenzy"
	//case UPA_SPARKSHOT:
	//	name = "Sparkshot"
	//case UPA_ANTI_AIR:
	//	name = "Anti-Air"
	default:
		name = "UNNAMED_ABILITY_" + strconv.Itoa(int(us.code))
	}
	if us.value > 0 {
		return fmt.Sprintf("%s %d", name, us.value)
	} else {
		return name
	}
}

type unitPassiveAbilityCode uint8

const (
	UPA_HASTE unitPassiveAbilityCode = iota
	UPA_HEALING
	UPA_FRENZY
	UPA_SPARKSHOT
	UPA_FLYING
	UPA_ANTI_AIR
	UPA_READINESS
	UPA_OVERPOWER
	UPA_INVISIBLE
	UPA_OBLITERATE
	UPA_COUNT
)
