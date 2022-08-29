package main

import (
	"log"
	"net"

	"github.com/MonteCarloClub/KBD/frame"
	api "github.com/MonteCarloClub/KBD/kitex_gen/api/kanbandatabase"
	"github.com/cloudwego/kitex/server"
)

func main() {
	//w, _ := os.OpenFile(path.Join("/", constant.DataDir, constant.LogFile), os.O_WRONLY|os.O_CREATE, 0755)
	//klog.SetOutput(w)
	frame.Init()
	svr := api.NewServer(new(KanBanDatabaseImpl), server.WithServiceAddr(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8288}))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
