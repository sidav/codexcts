package main

type aiPlayerController struct {
	controlsPlayer *player
}

func (ai *aiPlayerController) act(g *game) {
	switch g.currentPhase {
	case PHASE_MAIN:
	case PHASE_CODEX:
	}
}
