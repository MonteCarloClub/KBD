package main

import (
	"sync"

	"github.com/MonteCarloClub/KBD/config"
	"github.com/MonteCarloClub/KBD/log"
	"github.com/MonteCarloClub/KBD/rpcserver"
	"github.com/MonteCarloClub/KBD/signal"
)

var (
	cfg       *config.Config
	rpcServer *rpcserver.RpcServer
)

func main() {
	regs := []register{
		{"GetAccountData", GetAccountData},
		{"GetBlockData", GetBlockData},
	}
	newServer(1233, regs)
	wg := sync.WaitGroup{}
	interrupt := signal.InterruptListener()
	defer func() {
		log.AcbcLog.Infof("Gracefully shutting down the server...")
		if !cfg.DisableRPC {
			rpcServer.Stop()
		}
		//server.WaitForShutdown()
		log.SrvrLog.Infof("Server shutdown complete")
	}()
	if !cfg.DisableRPC {
		wg.Add(1)
		rpcServer.Start()
	}
	<-interrupt
	return
}
