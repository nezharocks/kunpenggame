package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	core "../core"
)

var (
	name = *flag.String("name", "吕红利", "set game server name")
	port = *flag.Int("port", 2019, "set game server port")
)

func main() {
	address := fmt.Sprintf(":%v", port)
	log.Printf("game server %q is starting ...", name)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		// todo shutdown gracefully until all connections are closed.
		log.Fatal(err)
	}
	log.Printf("game server %q is listenning on %v\n", name, address)

	gameImpl := core.NewGameImpl(name)
	gameService := core.NewGameService(gameImpl, ln)
	gameService.Serve()
}
