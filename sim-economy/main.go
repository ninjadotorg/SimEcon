package main

import (
	"log"

	"github.com/ninjadotorg/SimEcon/sim-economy/economy"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	economy.Run()
}
