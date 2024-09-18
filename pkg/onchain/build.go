package onchain

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/simplechain-org/client"
	"github.com/simplechain-org/client/common"
	"github.com/simplechain-org/client/core/types"
	"github.com/simplechain-org/client/crypto"
)

type Client interface {
	Close()
	ChainID(ctx context.Context) (*big.Int, error)
	PendingNonceAt(ctx context.Context, address common.Address) (uint64, error)
	EstimateGas(ctx context.Context, msg client.CallMsg) (uint64, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

func GetAddress(privateKey *ecdsa.PrivateKey) common.Address {
	return crypto.PubkeyToAddress(privateKey.PublicKey)
}
func BuildTx(chainClient Client, privateKey *ecdsa.PrivateKey, chainID *big.Int, nonce uint64, to common.Address, data []byte) (*types.Transaction, error) {
	accountAddress := GetAddress(privateKey)
	msg := client.CallMsg{From: accountAddress, To: &to, Value: big.NewInt(0), Data: data}
	gasLimit, err := chainClient.EstimateGas(context.TODO(), msg)
	if err != nil {
		return nil, err
	}
	tx := types.NewTransaction(nonce,
		to,
		big.NewInt(0), //amount=0
		gasLimit,
		big.NewInt(0), //gasPrice=0
		data,
	)
	return types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
}
