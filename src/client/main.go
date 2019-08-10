package main

import (
	"flag"

	core "../core"
)

var (
	id   = *flag.Int("id", 100, "set team id")
	name = *flag.String("name", "daolaji", "set team name")
	ip   = *flag.String("ip", "127.0.0.1", "set game server ip")
	port = *flag.Int("port", 2019, "set game server port")
)

func main() {
	team := core.NewTeamImpl(name)
	gameAgent := core.NewGameAgentImpl(id, name, ip, port)
	teamService := core.NewTeamService(team, gameAgent)
	teamService.Start()
}
