package main

func gameLoop(g *game) {
	io.renderGame(g, 0, pc)
	for !pc.exitGame {
		g.performCurrentPhase()
		for !pc.exitGame && g.currentPlayer == pc.controlsPlayer && !pc.phaseEnded {
			io.renderGame(g, g.currentPlayerNumber, pc)
			pc.act(g)
		}
		if g.currentPlayer == aiPc.controlsPlayer {
			aiPc.act(g)
		}
		pc.phaseEnded = false
		g.endCurrentPhase()
	}
}
