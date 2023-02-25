package main

func gameLoop(g *game) {
	for !exitGame {
		g.performCurrentPhase()
		pc := g.playersControllers[g.currentPlayerNumber]
		for !exitGame {
			pc.act(g)
			if pc.phaseEnded() {
				break
			}
		}
		g.endCurrentPhase()
	}
}
