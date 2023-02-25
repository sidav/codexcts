package main

const (
	PLAYERZONE_OTHER = iota
	PLAYERZONE_PATROL
	PLAYERZONE_TECH_BUILDINGS
	PLAYERZONE_ADDON_BUILDING
	PLAYERZONE_MAIN_BASE
)

type playerZoneCoords struct {
	zone        int
	indexInZone int
}
