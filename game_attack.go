package main

func (g *game) getAttackableCoordsForUnit(u *unit, owner *player) []*coords {
	enemy := g.getEnemyForPlayer(owner)
	if enemy.patrolZone[0] != nil {
		return []*coords{{PLAYERZONE_PATROL, 0}}
	}
	list := make([]*coords, 0)
	for i, p := range enemy.patrolZone {
		if p != nil {
			list = append(list, &coords{PLAYERZONE_PATROL, i})
		}
	}
	if len(list) == 0 { // patrol zone empty, adding everything the player has
		list = append(list, &coords{PLAYERZONE_MAIN_BASE, 0})
		for i, t := range enemy.techBuildings {
			if t != nil {
				list = append(list, &coords{PLAYERZONE_TECH_BUILDINGS, i})
			}
		}
		if enemy.addonBuilding != nil {
			list = append(list, &coords{PLAYERZONE_ADDON_BUILDING, 0})
		}
	}
	return list
}
