package household

import (
	"log"

	"github.com/0xroc/economy-simulation/state"
)

func Greedy(s state.State) {
	log.Println("i'm a greedy household")
}

func Modest(s state.State) {
	log.Println("i'm a modest household")
}
