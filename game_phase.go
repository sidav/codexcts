package main

const (
	PHASE_APPLY_TECH = iota
	PHASE_READY
	PHASE_UPKEEP
	PHASE_MAIN
	PHASE_DISCARD
	PHASE_CODEX
	TOTAL_PHASES
)

func (g *game) getCurrentPhaseName() string {
	switch g.currentPhase {
	case PHASE_APPLY_TECH:
		return "Apply tech"
	case PHASE_READY:
		return "Ready"
	case PHASE_UPKEEP:
		return "Upkeep"
	case PHASE_MAIN:
		return "Main"
	case PHASE_DISCARD:
		return "Discard"
	case PHASE_CODEX:
		return "Select tech"
	}
	panic("No phase name")
}
