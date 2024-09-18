package onchain

import (
	"container/heap"
	"math/big"

	"github.com/zeromicro/go-zero/core/logx"
)

func (m *Manager) HandleLastBlock(lastBlock string) {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	go func() {
		lastBlockNumber := big.NewInt(0)
		lastBlockNumber.SetString(lastBlock, 10)
		current := big.NewInt(0)
		one := big.NewInt(1)
		for {
			select {
			case number := <-m.handleNumbers:
				m.Logger.Infow("HandleLastBlock", logx.Field("number", number))
				heap.Push(&pq, number)
			inner:
				for {
					if pq.Len() <= 0 {
						break inner
					}
					item, ok := heap.Pop(&pq).(*big.Int)
					if !ok {
						break inner
					}
					if item.Cmp(lastBlockNumber) == 0 {
						break inner
					}
					current.Set(lastBlockNumber)
					current.Add(current, one)

					m.Logger.Infow("checkBlock", logx.Field("current", current.String()), logx.Field("item", item.String()))
					if item.Cmp(current) == 0 {
						m.Logger.Infow("expected", logx.Field("number", item.String()))
						lastBlockNumber.Set(item)
						//是预期的区块
						if m.blockUpdater != nil {
							m.blockUpdater.OnBlock(current)
						}
					} else {
						if item.Cmp(lastBlockNumber) > 0 {
							//不是预期的并且要大于已经处理过的，先放回去
							heap.Push(&pq, item)
							m.Logger.Infow("unexpected", logx.Field("number", item.String()))
						}
						break inner
					}
				}
			}
		}
	}()
}
