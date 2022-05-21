package types

import (
	"fmt"
	"math/big"

	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/state"
)

type BlockProcessor interface {
	Process(*Block) (state.Logs, error)
}

const bloomLength = 256

type Bloom [bloomLength]byte

func BytesToBloom(b []byte) Bloom {
	var bloom Bloom
	bloom.SetBytes(b)
	return bloom
}

func (b *Bloom) SetBytes(d []byte) {
	if len(b) < len(d) {
		panic(fmt.Sprintf("bloom bytes too big %d %d", len(b), len(d)))
	}

	copy(b[bloomLength-len(d):], d)
}

func (b Bloom) Big() *big.Int {
	return common.Bytes2Big(b[:])
}

func (b Bloom) Bytes() []byte {
	return b[:]
}
