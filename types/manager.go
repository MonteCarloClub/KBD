package types

import (
	"github.com/MonteCarloClub/KBD/accounts"
	"github.com/MonteCarloClub/KBD/chain_manager"
	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/event"
	"github.com/MonteCarloClub/KBD/kbpool"
)

type Backend interface {
	AccountManager() *accounts.Manager
	BlockProcessor() *BlockProcessor
	ChainManager() *chain_manager.Manager
	TxPool() *kbpool.TxPool
	BlockDb() common.Database
	StateDb() common.Database
	EventMux() *event.TypeMux
}
