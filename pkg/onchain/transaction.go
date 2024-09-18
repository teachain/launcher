package onchain

import (
	"github.com/simplechain-org/client/common"
	"github.com/zeromicro/go-zero/core/logx"
	"math/big"
	"os"
	"time"
)

func (m *Manager) TransactionFinished(hash common.Hash, blockNumber *big.Int) error {
	m.onChainTxCh <- hash.String()
	return nil
}
func (m *Manager) TransactionSent() {
	go func() {
		sentCount := 0
		sentTx := make(map[string]struct{})
		var onChainSuccess int = 0
		for {
			select {
			case tx := <-m.txCh:
				sentTx[tx] = struct{}{}
				sentCount++
				elapsed := time.Now().Sub(m.startedAt)
				seconds := (int)(elapsed.Seconds())
				if seconds > 0 {
					tps := sentCount / seconds
					logx.Infow("send", logx.Field("tps", tps), logx.Field("sentCount", sentCount), logx.Field("elapsed", elapsed.String()))
				}

			case tx := <-m.onChainTxCh:
				if _, ok := sentTx[tx]; ok {
					onChainSuccess++
					if m.chainConfig.TxCount == onChainSuccess && !m.chainConfig.NoLimit {
						elapsed := time.Now().Sub(m.startedAt)
						seconds := (int)(elapsed.Seconds())
						if seconds > 0 {
							tps := onChainSuccess / seconds
							logx.Infow("finally", logx.Field("tps", tps), logx.Field("elapsed", elapsed.String()))
							//跑完了，程序退出
							os.Exit(0)
						}
					}
					elapsed := time.Now().Sub(m.startedAt)
					seconds := (int)(elapsed.Seconds())
					if seconds > 0 {
						tps := onChainSuccess / seconds
						logx.Infow("runtime", logx.Field("tps", tps), logx.Field("onChainSuccess", onChainSuccess), logx.Field("elapsed", elapsed.String()))
					}
					delete(sentTx, tx)
				}
			}
		}
	}()
}
