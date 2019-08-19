package main

import (
	"fmt"
	"log"
	"time"

	core ".."
)

func main() {
	start := time.Now()
	game := core.NewFirstGame(core.Map1, 4, 20, 20)
	err := game.Init()
	if err != nil {
		log.Println(err)
	}
	guest := core.NewTeamImpl("daolaji")
	guest.SetTeamID(game.NewTeamID())
	host := core.NewTeamImpl("ai")
	host.SetTeamID(game.NewTeamID())
	battle := game.NewBattle(guest, host)
	battle.Run()
	fmt.Println("time spent:", time.Since(start))
}
