package main

import (
	"fmt"
	"strconv"
)

func (g *game) putSpellInEffect(owner *player, s *spell, targetCoords *playerZoneCoords) {
	switch s.effectCode {
	case SPELL_DEAL_DAMAGE:
		g.dealDamageByCoords(s.effectValue1, targetCoords)
	case SPELL_ADD_RUNES:
		targUnt := targetCoords.player.getUnitByCoords(targetCoords)
		r1Text := strconv.Itoa(s.effectValue1)
		if s.effectValue1 >= 0 {
			r1Text = "+" + r1Text
		}
		r2Text := strconv.Itoa(s.effectValue1)
		if s.effectValue1 >= 0 {
			r2Text = "+" + r2Text
		}
		g.messageForPlayer += fmt.Sprintf("%s receives %s/%s rune.", targUnt.getName(), r1Text, r2Text)
		targUnt.receiveRune(s.effectValue1, s.effectValue2)
	}
	g.removeDeadUnits()
}

func (g *game) getTargetableCoordsForSpell(owner *player, s *spell) []*playerZoneCoords {
	list := make([]*playerZoneCoords, 0)
	for _, currPlayer := range g.players {
		if currPlayer == owner && !s.canTargetOwner {
			continue
		}
		if currPlayer != owner && !s.canTargetEnemy {
			continue
		}
		for _, z := range s.targetableZones {
			switch z {
			case PLAYERZONE_OTHER:
				for i, unt := range currPlayer.otherZone {
					if g.canSpellTargetUnit(s, unt) {
						list = append(list, &playerZoneCoords{
							player:      currPlayer,
							zone:        PLAYERZONE_OTHER,
							indexInZone: i,
						})
					}
				}
			case PLAYERZONE_PATROL:
				for i, patroller := range currPlayer.patrolZone {
					if patroller != nil && g.canSpellTargetUnit(s, patroller) {
						list = append(list, &playerZoneCoords{
							player:      currPlayer,
							zone:        PLAYERZONE_PATROL,
							indexInZone: i,
						})
					}
				}
			}
		}
	}
	return list
}

func (g *game) canSpellTargetUnit(s *spell, u *unit) bool {
	return s.canTargetHeroes && u.isHero() || s.canTargetUnits && !u.isHero()
}
