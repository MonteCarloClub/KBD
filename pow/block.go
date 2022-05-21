package pow

import (
	"math/big"

	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/types"
)

type Block interface {
	Difficulty() *big.Int
	HashNoNonce() common.Hash
	Nonce() uint64
	MixDigest() common.Hash
	NumberU64() uint64
}

type ChainManager interface {
	GetBlockByNumber(uint64) *types.Block
	CurrentBlock() *types.Block
}
