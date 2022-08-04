package crypto

// encryption algorithm interface
type EA interface {
	GenerateKeyPair() ([]byte, []byte)
	GeneratePubKey(seckey []byte) ([]byte, error)
	RecoverPubkey(msg []byte, sig []byte)
	Sign(msg []byte, seckey []byte) ([]byte, error)
}
