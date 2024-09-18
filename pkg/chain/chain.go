package chain

import (
	"context"
	"github.com/simplechain-org/client/common"
	"github.com/simplechain-org/client/common/hexutil"
	"github.com/simplechain-org/client/rpc"
)

func NewClient(endpoint string) (Client, error) {
	c, err := rpc.DialContext(context.Background(), endpoint)
	if err != nil {
		return nil, err
	}
	return &client{client: c, Endpoint: endpoint}, nil
}

type client struct {
	Endpoint string
	client   *rpc.Client
}

func (c *client) TxPoolStatus() (*TxPoolStatusResponse, error) {
	result := new(TxPoolStatusResponse)
	err := c.client.Call(result, "txpool_status")
	return result, err
}
func (c *client) TransactionCountLatest(ctx context.Context, account common.Address) (uint64, error) {
	var result hexutil.Uint64
	err := c.client.CallContext(ctx, &result, "eth_getTransactionCount", account, "latest")
	return uint64(result), err
}
