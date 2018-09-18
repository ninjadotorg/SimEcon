package main

import (
	"log"

	"github.com/ninjadotorg/SimEcon/sim-production/production"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	production.Run()
}
