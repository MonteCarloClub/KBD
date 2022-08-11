package main

import (
	"log"
	"net"

	api "github.com/MonteCarloClub/KBD/kitex_gen/api/kanbandatabase"
	"github.com/cloudwego/kitex/server"
)

func main() {
	svr := api.NewServer(new(KanBanDatabaseImpl), server.WithServiceAddr(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8288}))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
