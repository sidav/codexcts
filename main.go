package main

import (
	"codexcts/lib/random"
	"codexcts/lib/random/pcgrandom"
	"fmt"
	"log"
	"os"
)

var (
	rnd  random.PRNG
	pc   *playerController
	aiPc *aiPlayerController
)

func main() {
	onInit()
	defer onClose()

	//f, err := os.OpenFile(fmt.Sprintf("debug_output%s.log", time.Now().Format("2006_01_02_15_04_05")),
	//	os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	f, err := os.OpenFile("debug_output.log",
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(fmt.Sprintf("Error opening file: %v", err))
	}
	defer f.Close()
	log.SetOutput(f)

	rnd = pcgrandom.New(-1)
	g := &game{}
	g.initGame()
	pc = &playerController{
		controlsPlayer:              g.players[0],
		currentMode:                 PCMODE_NONE,
		currentSelectedCardFromHand: nil,
		currentSelectedUnit:         nil,
	}
	aiPc = &aiPlayerController{
		controlsPlayer: g.players[1],
	}
	gameLoop(g)
}
