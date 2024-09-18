package onchain

import (
	"crypto/ecdsa"
	"math/big"
)

type Job interface {
	Do(privateKey *ecdsa.PrivateKey, chainID *big.Int, nonce uint64) error
}
