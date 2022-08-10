package main

import (
	api "github.com/MonteCarloClub/KBD/kitex_gen/api/kanbandatabase"
	"log"
)

func main() {
	svr := api.NewServer(new(KanBanDatabaseImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
