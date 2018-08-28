package main

import (
	"flag"
	"log"

	"github.com/0xroc/economy-simulation/sim"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	file := flag.String("s", "", "simulation definition")
	flag.Parse()

	if e := sim.Start(*file); e != nil {
		log.Println(e)
	}
}
