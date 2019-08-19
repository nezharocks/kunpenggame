package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	core "../core"
)

var (
	id, port int
	name, ip string
)

func main() {
	args := os.Args[1:]
	if len(args) < 3 {
		id = *flag.Int("id", 100, "set team id")
		name = *flag.String("name", "daolaji", "set team name")
		ip = *flag.String("ip", "127.0.0.1", "set game server ip")
		port = *flag.Int("port", 2019, "set game server port")
	} else {
		id, _ = strconv.Atoi(args[0])
		ip = args[1]
		port, _ = strconv.Atoi(args[2])
		name = "daolaji"
	}

	team := core.NewTeamImpl(name)
	team.SetTeamID(id)
	gameAgent := core.NewGameAgentImpl(id, name, ip, port)
	teamService := core.NewTeamService(team, gameAgent)

	start := time.Now()
	teamService.Start()
	log.Println("time spent:", time.Since(start))
}
