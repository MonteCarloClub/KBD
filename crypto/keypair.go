package crypto

import (
	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/crypto/secp256k1"
	"strings"
)

type KeyPair struct {
	PrivateKey []byte
	PublicKey  []byte
	address    []byte
	mnemonic   string
	// The associated account
	// account *StateObject
}

func GenerateNewKeyPair() *KeyPair {
	_, prv := secp256k1.GenerateKeyPair()
	keyPair, _ := NewKeyPairFromSec(prv) // swallow error, this one cannot err
	return keyPair
}

func NewKeyPairFromSec(seckey []byte) (*KeyPair, error) {
	pubkey, err := secp256k1.GeneratePubKey(seckey)
	if err != nil {
		return nil, err
	}

	return &KeyPair{PrivateKey: seckey, PublicKey: pubkey}, nil
}

func (k *KeyPair) Address() []byte {
	if k.address == nil {
		k.address = Sha3(k.PublicKey[1:])[12:]
	}
	return k.address
}

func (k *KeyPair) Mnemonic() string {
	if k.mnemonic == "" {
		k.mnemonic = strings.Join(MnemonicEncode(common.Bytes2Hex(k.PrivateKey)), " ")
	}
	return k.mnemonic
}

func (k *KeyPair) AsStrings() (string, string, string, string) {
	return k.Mnemonic(), common.Bytes2Hex(k.Address()), common.Bytes2Hex(k.PrivateKey), common.Bytes2Hex(k.PublicKey)
}
