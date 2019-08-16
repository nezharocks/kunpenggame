package main

import (
	"fmt"
	"log"

	core ".."
)

func main() {
	fmt.Println("hi")
	j := core.NewFirstGame(core.Map1, 4, 20, 20)
	err := j.Init()
	if err != nil {
		log.Println(err)
	}
	t1 := core.NewTeamImpl("ai_team")
	t1.SetTeamID(j.NewTeamID())
	t2 := core.NewTeamImpl("ai_team")
	t2.SetTeamID(j.NewTeamID())
	// ts := []core.TeamBattle{t1, t2}

	battle := j.NewBattle(t1, t2)
	battle.Run()
	// legs := battle.GetLegs()
	// for _, leg := range legs {
	// 	legStart := leg.Start()
	// 	fmt.Printf("%+v", legStart.Message())
	// 	t1.LegStart(legStart)
	// 	t2.LegStart(legStart)

	// 	var (
	// 		err   error
	// 		round *core.Round
	// 		// action *core.Action
	// 	)
	// 	beatTeam := t1
	// 	thinkTeam := t2
	// 	semiRoundNum := leg.RoundNum

	// 	for i := 0; i < semiRoundNum; i++ {
	// 		// think activities
	// 		round = leg.Round(i, thinkTeam.GetTeamID(), core.ThinkMode)
	// 		err = thinkTeam.Round(round)
	// 		if err != nil {
	// 			// todo
	// 			log.Println(err)
	// 		}
	// 	wait_think:
	// 		select {
	// 		case thinkAction := <-thinkTeam.GetActionCh():
	// 			fmt.Println(thinkAction.Message())
	// 			leg.Action(&thinkAction)
	// 		case <-time.After(time.Millisecond * 800):
	// 			break wait_think
	// 		}

	// 		// beat activities
	// 		round = leg.Round(i, beatTeam.GetTeamID(), core.BeatMode)
	// 		err = beatTeam.Round(round)
	// 		if err != nil {
	// 			// todo
	// 			log.Println(err)
	// 		}
	// 	wait_beat:
	// 		select {
	// 		case beatAction := <-beatTeam.GetActionCh():
	// 			fmt.Println(beatAction.Message())
	// 			leg.Action(&beatAction)
	// 		case <-time.After(time.Millisecond * 800):
	// 			break wait_beat
	// 		}
	// 	}

	// 	legEnd := leg.End()
	// 	fmt.Printf("%+v", legEnd.Message())
	// 	t1.LegEnd(legEnd)
	// 	t2.LegEnd(legEnd)
	// }
	// fmt.Printf("%+v", legs[0].JSON())

}
