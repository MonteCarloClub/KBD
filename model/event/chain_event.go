package event

import (
	"math/big"

	"github.com/MonteCarloClub/KBD/model/state"

	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/types"
)

// TxPreEvent is posted when a transaction enters the transaction pool.
type TxPreEvent struct{ Tx *types.Transaction }

// TxPostEvent is posted when a transaction has been processed.
type TxPostEvent struct{ Tx *types.Transaction }

// NewBlockEvent is posted when a block has been imported.
type NewBlockEvent struct{ Block *types.Block }

// NewMinedBlockEvent is posted when a block has been imported.
type NewMinedBlockEvent struct{ Block *types.Block }

// ChainSplit is posted when a new head is detected
type ChainSplitEvent struct {
	Block *types.Block
	Logs  state.Logs
}

type ChainEvent struct {
	Block *types.Block
	Hash  common.Hash
	Logs  state.Logs
}

type ChainSideEvent struct {
	Block *types.Block
	Logs  state.Logs
}

type PendingBlockEvent struct {
	Block *types.Block
	Logs  state.Logs
}

type ChainUncleEvent struct {
	Block *types.Block
}

type ChainHeadEvent struct{ Block *types.Block }

type GasPriceChanged struct{ Price *big.Int }

// Mining operation events
type StartMining struct{}
type TopMining struct{}
