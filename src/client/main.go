package main

import (
	"flag"
	"fmt"
	"log"

	core "../core"
)

var (
	id   = *flag.Int("id", 100, "set team id")
	name = *flag.String("name", "daolaji", "set team name")
	ip   = *flag.String("ip", "127.0.0.1", "set game server ip")
	port = *flag.Int("port", 2019, "set game server port")
)

func main() {
	address := fmt.Sprintf("%v:%v", ip, port)
	log.Printf("team %q is connecting to game server on %v ...", name, address)
	team := core.NewTeamSimple(name)
	gameAgent := core.NewGameAgentImpl(id, name, ip, port)
	err := gameAgent.Connect()
	if err != nil {
		log.Println(err)
		return
	}
	teamService := core.NewTeamService(team, gameAgent)
	teamService.Start()
}
