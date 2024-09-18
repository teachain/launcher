package onchain

import (
	"context"
	"github.com/simplechain-org/client/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
	"math/big"
	"time"
)

// UpdateOnChainStatus
// 1.定时
// 2.有新区块产生
func (m *Manager) UpdateOnChainStatus(lastBlock string) {
	client, err := ethclient.DialContext(context.Background(), m.chainConfig.WsEndpoint)
	if err != nil {
		m.Logger.Errorw("chain client DialContext", logx.Field("error", err.Error()))
		return
	}
	numbers := make(chan *big.Int, 100)
	blockBrowser := NewBlockBrowser(client, time.Second*5, func(number *big.Int) {
		select {
		case numbers <- number:
		default:
			m.Logger.Error("numbers is busy")
		}
	})
	err = blockBrowser.Start()
	if err != nil {
		m.Logger.Errorw("blockBrowser start", logx.Field("error", err.Error()))
		return
	}
	go func() {
		interval := time.Minute
		timer := time.NewTimer(interval)
		defer timer.Stop()

		maxNumber := big.NewInt(0)
		start := big.NewInt(0)

		start.SetString(lastBlock, 10)

		one := big.NewInt(1)

		blockNumber, err := client.BlockNumber(context.Background())
		if err != nil {
			m.Logger.Errorw("BlockNumber", logx.Field("error", err.Error()))
			return
		}
		maxNumber.SetUint64(blockNumber)
		m.HandleLastBlock(lastBlock)
		for {
			for maxNumber.Cmp(start) >= 0 {
				m.dispatcher.Enqueue(&Task{
					number:            big.NewInt(0).Set(start),
					endpoint:          m.chainConfig.HttpEndpoint,
					numbers:           m.handleNumbers,
					handleTransaction: m.TransactionFinished,
				})
				start.Add(start, one)
			}
			select {
			case <-timer.C:
				//要么时间到
				//查一下区块高度
				blockNumber, err := client.BlockNumber(context.Background())
				if err != nil {
					m.Logger.Errorw("BlockNumber", logx.Field("error", err.Error()))
					time.Sleep(time.Second * 5)
					timer.Reset(interval)
					continue
				}
				number := big.NewInt(0).SetUint64(blockNumber)
				if number.Cmp(maxNumber) > 0 {
					maxNumber.Set(number)
				}
				timer.Reset(interval)
			case number := <-numbers:
				if number.Cmp(maxNumber) > 0 {
					maxNumber.Set(number)
				}
				//要么有新区块到来
				timer.Stop()
				timer.Reset(interval)
				m.Logger.Debugw("New block is coming", logx.Field("number", number))
			}
		}
	}()
}
