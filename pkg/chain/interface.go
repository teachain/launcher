package chain

import (
	"context"
	"github.com/simplechain-org/client/common"
)

type Client interface {
	TxPoolStatus() (*TxPoolStatusResponse, error)
	TransactionCountLatest(ctx context.Context, account common.Address) (uint64, error)
}
