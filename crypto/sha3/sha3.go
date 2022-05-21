package sha3

import (
	"hash"

	"golang.org/x/crypto/sha3"
)

func Sha3(data []byte) []byte {
	d := sha3.New256()
	d.Write(data)

	return d.Sum(nil)
}
func NewKeccak256() hash.Hash {
	return sha3.New256()
}
