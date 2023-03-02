package main

import (
	"codexcts/lib/random"
	"codexcts/lib/random/pcgrandom"
	"fmt"
	golangIo "io"
	"log"
	"os"
)

var (
	exitGame bool
	rnd      random.PRNG
)

func main() {
	onInit()
	defer onClose()

	if len(os.Args) > 1 && (os.Args[1] == "log" || os.Args[1] == "aivsai") {
		f, err := os.OpenFile("debug_output.log",
			os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			panic(fmt.Sprintf("Error opening file: %v", err))
		}
		defer f.Close()
		log.SetOutput(f)
	} else {
		log.SetOutput(golangIo.Discard)
	}

	rnd = pcgrandom.New(-1)
	g := &game{}
	g.initGame()

	if len(os.Args) > 1 && os.Args[1] == "aivsai" {
		setupGame(g, true)
	} else {
		setupGame(g, false)
	}

	// debug
	// g.players[0].gold += 200
	//g.players[0].otherZone = append(g.players[0].otherZone, &unit{
	//	card: getCardByName("Eggship"),
	//})
	//g.players[0].otherZone = append(g.players[0].otherZone, &unit{
	//	card: getCardByName("Eggship"),
	//})
	//g.players[0].otherZone = append(g.players[0].otherZone, &unit{
	//	card: getCardByName("Leaping Lizard"),
	//})
	//g.players[1].patrolZone[0] = &unit{
	//	card: getCardByName("Eggship"),
	//}

	gameLoop(g)
}

func setupGame(g *game, aivsai bool) {
	if rnd.OneChanceFrom(2) {
		g.players[0].commandZone[0] = heroCardsDb[1]
		g.players[1].commandZone[0] = heroCardsDb[0]
	} else {
		g.players[0].commandZone[0] = heroCardsDb[0]
		g.players[1].commandZone[0] = heroCardsDb[1]
	}
	g.initPlayerCodices()

	if aivsai {
		g.players[0].name = "AI 1"
		g.playersControllers = append(g.playersControllers, &aiPlayerController{
			controlsPlayer: g.players[0],
		})
		g.players[1].name = "AI 2"
		g.playersControllers = append(g.playersControllers, &aiPlayerController{
			controlsPlayer: g.players[1],
		})
		return
	}

	playerFirst := rnd.OneChanceFrom(2)
	if playerFirst {
		g.messageForPlayer = "You play first. \n"
		g.messageForPlayer += fmt.Sprintf("Your starting hero is %s. \n ", g.players[0].commandZone[0].getName())
		g.players[0].name = "Player"
		g.playersControllers = append(g.playersControllers, &playerController{
			g:                           g,
			controlsPlayer:              g.players[0],
			currentMode:                 PCMODE_NONE,
			currentSelectedCardFromHand: nil,
			currentSelectedUnit:         nil,
		})

		g.players[1].name = "AI"
		g.playersControllers = append(g.playersControllers, &aiPlayerController{
			controlsPlayer: g.players[1],
		})
	} else {
		g.players[0].name = "AI"
		g.playersControllers = append(g.playersControllers, &aiPlayerController{
			controlsPlayer: g.players[0],
		})

		g.messageForPlayer = "You play second. \n "
		g.messageForPlayer += fmt.Sprintf("Your starting hero is %s. \n ", g.players[1].commandZone[0].getName())
		g.players[1].name = "Player"
		g.playersControllers = append(g.playersControllers, &playerController{
			g:                           g,
			controlsPlayer:              g.players[1],
			currentMode:                 PCMODE_NONE,
			currentSelectedCardFromHand: nil,
			currentSelectedUnit:         nil,
		})
	}
	g.showMessageToAllPlayers("GAME START")
}
