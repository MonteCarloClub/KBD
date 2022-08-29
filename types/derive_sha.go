package types

import (
	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/model/kdb"
	"github.com/MonteCarloClub/KBD/model/trie"
	"github.com/MonteCarloClub/KBD/rlp"
)

type DerivableList interface {
	Len() int
	GetRlp(i int) []byte
}

func DeriveSha(list DerivableList) common.Hash {
	db, _ := kdb.NewMemDatabase()
	trie := trie.New(nil, db)
	for i := 0; i < list.Len(); i++ {
		key, _ := rlp.EncodeToBytes(uint(i))
		trie.Update(key, list.GetRlp(i))
	}

	return common.BytesToHash(trie.Root())
}
