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
		g.players[0].name = "Player 1"
		g.playersControllers = append(g.playersControllers, &aiPlayerController{
			controlsPlayer: g.players[0],
		})
	} else {
		g.players[0].name = "Player"
		g.playersControllers = append(g.playersControllers, &playerController{
			controlsPlayer:              g.players[0],
			currentMode:                 PCMODE_NONE,
			currentSelectedCardFromHand: nil,
			currentSelectedUnit:         nil,
		})
	}

	g.playersControllers = append(g.playersControllers, &aiPlayerController{
		controlsPlayer: g.players[1],
	})
	g.players[1].name = "AI - Player 2"

	// debug units
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
