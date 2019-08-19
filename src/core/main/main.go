package main

import (
	"fmt"
	"log"

	core ".."
)

func main() {
	fmt.Println("hi")
	game := core.NewFirstGame(core.Map1, 4, 20, 20)
	err := game.Init()
	if err != nil {
		log.Println(err)
	}
	guest := core.NewTeamImpl("daolaji")
	guest.SetTeamID(game.NewTeamID())
	host := core.NewTeamImpl("ai")
	host.SetTeamID(game.NewTeamID())
	// ts := []core.TeamBattle{t1, t2}

	battle := game.NewBattle(guest, host)
	battle.Run()
}
