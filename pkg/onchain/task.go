package onchain

import (
	"context"
	"github.com/simplechain-org/client/common"
	"github.com/simplechain-org/client/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
	"math/big"
	"time"
)

type Task struct {
	number            *big.Int //区块高度
	endpoint          string
	handleTransaction func(hash common.Hash, blockNumber *big.Int) error
	numbers           chan *big.Int
}

func (t *Task) Do() {
	for {
		logx.Infow("Task do number", logx.Field("number", t.number.String()))
		client, err := ethclient.DialContext(context.Background(), t.endpoint)
		if err != nil {
			logx.Errorw("DialContext", logx.Field("error", err.Error()))
			time.Sleep(time.Second * 5)
			continue
		}
		//有可能
		block, err := client.BlockByNumber(context.Background(), t.number)
		if err != nil {
			client.Close()
			logx.Errorw("BlockByNumber", logx.Field("error", err.Error()), logx.Field("number", t.number))
			time.Sleep(time.Second * 5)
			continue
		}
		transactions := block.Transactions()
	loop:
		for _, tx := range transactions {
			if t.handleTransaction != nil {
				err := t.handleTransaction(tx.Hash(), t.number)
				if err != nil {
					logx.Errorw("UpdateStatusByTxHash", logx.Field("error", err.Error()))
					time.Sleep(time.Second * 5)
					continue loop
				}
			}
		}
		logx.Infow("Task join numbers", logx.Field("number", t.number.String()))
		t.numbers <- t.number
		client.Close()
		return
	}
}
