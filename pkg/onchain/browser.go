package onchain

import (
	"context"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/simplechain-org/client/core/types"
	"github.com/simplechain-org/client/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
)

type BlockBrowser struct {
	client     *ethclient.Client
	timeout    time.Duration
	stop       chan struct{}
	isStopped  int32
	onNewBlock func(Number *big.Int)
}

func NewBlockBrowser(client *ethclient.Client, timeout time.Duration, onNewBlock func(Number *big.Int)) *BlockBrowser {
	b := &BlockBrowser{
		client:     client,
		timeout:    timeout,
		stop:       make(chan struct{}),
		isStopped:  0,
		onNewBlock: onNewBlock,
	}
	return b
}

func (b *BlockBrowser) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), b.timeout)
	defer cancel()
	ch := make(chan *types.Header, 10)
	subscription, err := b.client.SubscribeNewHead(ctx, ch)
	if err != nil {
		logx.Errorw("SubscribeNewHead", logx.Field("error", err.Error()))
		return err
	}
	go func() {
		defer func() {
			logx.Info("BlockBrowser Start exited")
		}()
		for {
			select {
			case header := <-ch:

				logx.Infow("Receive new block",
					logx.Field("number", header.Number.String()),
					logx.Field("at", time.Now().Format(time.DateTime)))

				b.onNewBlock(header.Number)
			case err := <-subscription.Err():
				if err != nil {
					logx.Errorw("SubscribeNewHead run", logx.Field("error", err.Error()))
				loop:
					for {
						select {
						case <-b.stop:
							logx.Info("BlockBrowser Start subscribeNewHead exits")
							return
						default:
							//连接断开以后，需要重新进行订阅，虽然Client自身有重连操作，但不会自动重新订阅
							subscription, err = b.client.SubscribeNewHead(context.Background(), ch)
							if err != nil {
								logx.Errorw("SubscribeNewHead Resubscribe", logx.Field("error", err.Error()))
								time.Sleep(time.Second * 5)
							} else {
								break loop
							}
						}
					}
				}
			case <-b.stop:
				return
			}
		}
	}()
	return nil
}

func (b *BlockBrowser) Stop() {
	if atomic.CompareAndSwapInt32(&b.isStopped, 0, 1) {
		close(b.stop)
	}
}
