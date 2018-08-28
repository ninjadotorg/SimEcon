package main

import (
	"flag"
	"log"

	"github.com/0xroc/economy-simulation/economy"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	file := flag.String("f", "", "economy description file")
	flag.Parse()

	if e := economy.Run(*file); e != nil {
		log.Println(e)
	}
}
