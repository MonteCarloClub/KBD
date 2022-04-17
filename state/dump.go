/*
Copyright (c) 2022 Zhu Zunxiong <liuzunxiong@qq.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package state

import (
	"encoding/json"
	"fmt"

	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/common/hexutil"
	"github.com/MonteCarloClub/KBD/log"
	"github.com/MonteCarloClub/KBD/rlp"
	"github.com/MonteCarloClub/KBD/trie"
)

// DumpAccount represents an account in the state
type DumpAccount struct {
	Balance   string                 `json:"balance"`
	Nonce     uint64                 `json:"nonce"`
	Root      string                 `json:"root"`
	CodeHash  string                 `json:"codeHash"`
	Code      string                 `json:"code,omitempty"`
	Storage   map[common.Hash]string `json:"storage,omitempty"`
	Address   *common.Address        `json:"address,omitempty"` // Address only present in iterative (line-by-line) mode
	SecureKey hexutil.Bytes          `json:"key,omitempty"`     // If we don't have address, we can output the key

}

// Dump represents the full dump in a collected format, as one large map
type Dump struct {
	Root     string                         `json:"root"`
	Accounts map[common.Address]DumpAccount `json:"accounts"`
}

// iterativeDump is a 'collector'-implementation which dump output line-by-line iteratively
type iterativeDump json.Encoder

// Collector interface which the state trie calls during iteration
type collector interface {
	onRoot(common.Hash)
	onAccount(common.Address, DumpAccount)
}

func (self *Dump) onRoot(root common.Hash) {
	self.Root = fmt.Sprintf("%x", root)
}

func (self *Dump) onAccount(addr common.Address, account DumpAccount) {
	self.Accounts[addr] = account
}

func (self iterativeDump) onAccount(addr common.Address, account DumpAccount) {
	dumpAccount := &DumpAccount{
		Balance:   account.Balance,
		Nonce:     account.Nonce,
		Root:      account.Root,
		CodeHash:  account.CodeHash,
		Code:      account.Code,
		Storage:   account.Storage,
		SecureKey: account.SecureKey,
		Address:   nil,
	}
	if addr != (common.Address{}) {
		dumpAccount.Address = &addr
	}
	(*json.Encoder)(&self).Encode(dumpAccount)
}
func (self iterativeDump) onRoot(root common.Hash) {
	(*json.Encoder)(&self).Encode(struct {
		Root common.Hash `json:"root"`
	}{root})
}

func (self *StateDB) dump(c collector, excludeCode, excludeStorage, excludeMissingPreimages bool) {
	emptyAddress := (common.Address{})
	missingPreimages := 0
	c.onRoot(self.trie.Hash())
	it := trie.NewIterator(self.trie.NodeIterator(nil))
	for it.Next() {
		var data Account
		if err := rlp.DecodeBytes(it.Value, &data); err != nil {
			panic(err)
		}
		addr := common.BytesToAddress(self.trie.GetKey(it.Key))
		obj := newObject(nil, addr, data)
		account := DumpAccount{
			Balance:  data.Balance.String(),
			Nonce:    data.Nonce,
			Root:     common.Bytes2Hex(data.Root[:]),
			CodeHash: common.Bytes2Hex(data.CodeHash),
		}
		if emptyAddress == addr {
			// Preimage missing
			missingPreimages++
			if excludeMissingPreimages {
				continue
			}
			account.SecureKey = it.Key
		}
		if !excludeCode {
			account.Code = common.Bytes2Hex(obj.Code(self.db))
		}
		if !excludeStorage {
			account.Storage = make(map[common.Hash]string)
			storageIt := trie.NewIterator(obj.getTrie(self.db).NodeIterator(nil))
			for storageIt.Next() {
				_, content, _, err := rlp.Split(storageIt.Value)
				if err != nil {
					log.Error("Failed to decode the value returned by iterator", "error", err)
					continue
				}
				account.Storage[common.BytesToHash(self.trie.GetKey(storageIt.Key))] = common.Bytes2Hex(content)
			}
		}
		c.onAccount(addr, account)
	}
	if missingPreimages > 0 {
		log.Warn("Dump incomplete due to missing preimages", "missing", missingPreimages)
	}
}

// RawDump returns the entire state an a single large object
func (self *StateDB) RawDump(excludeCode, excludeStorage, excludeMissingPreimages bool) Dump {
	dump := &Dump{
		Accounts: make(map[common.Address]DumpAccount),
	}
	self.dump(dump, excludeCode, excludeStorage, excludeMissingPreimages)
	return *dump
}

// Dump returns a JSON string representing the entire state as a single json-object
func (self *StateDB) Dump(excludeCode, excludeStorage, excludeMissingPreimages bool) []byte {
	dump := self.RawDump(excludeCode, excludeStorage, excludeMissingPreimages)
	json, err := json.MarshalIndent(dump, "", "    ")
	if err != nil {
		fmt.Println("dump err", err)
	}
	return json
}

// IterativeDump dumps out accounts as json-objects, delimited by linebreaks on stdout
func (self *StateDB) IterativeDump(excludeCode, excludeStorage, excludeMissingPreimages bool, output *json.Encoder) {
	self.dump(iterativeDump(*output), excludeCode, excludeStorage, excludeMissingPreimages)
}
