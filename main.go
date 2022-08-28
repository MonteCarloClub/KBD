package main

import (
	"github.com/astaxie/beego/logs"
	"log"
	"net"

	api "github.com/MonteCarloClub/KBD/kitex_gen/api/kanbandatabase"
	"github.com/cloudwego/kitex/server"
)

func main() {
	logs.SetLogger(logs.AdapterFile, `{"filename":"project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)

	svr := api.NewServer(new(KanBanDatabaseImpl), server.WithServiceAddr(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8288}))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
