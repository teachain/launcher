package onchain

import (
	"context"
	"crypto/ecdsa"
	"github.com/simplechain-org/client/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
	"math/big"
	"strings"
	"time"
)

// Worker Worker结构体
// WorkerPool随机选取一个Worker，将Job发送给Worker去执行
type Worker struct {
	// 不需要带缓冲的任务队列
	jobQueue chan Job
	//退出标志
	stop chan struct{}

	privateKey *ecdsa.PrivateKey

	nonce uint64

	endpoint string

	chainId *big.Int
}

// NewWorker 创建一个新的Worker对象
func NewWorker(endpoint string, privateKey *ecdsa.PrivateKey) *Worker {
	return &Worker{
		jobQueue:   make(chan Job),
		stop:       make(chan struct{}),
		privateKey: privateKey,
		endpoint:   endpoint,
		chainId:    big.NewInt(0),
	}
}

func (w *Worker) GetNonce() (uint64, error) {
	client, err := ethclient.DialContext(context.Background(), w.endpoint)
	if err != nil {
		return 0, err
	}
	defer client.Close()

	accountAddress := GetAddress(w.privateKey)

	nonce, err := client.PendingNonceAt(context.TODO(), accountAddress)
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

// Start 启动一个Worker，来监听Job事件
// 执行完任务，需要将自己重新发送到WorkerPool
func (w *Worker) Start(workerPool *Dispatcher) {
	client, err := ethclient.DialContext(context.Background(), w.endpoint)
	if err != nil {
		logx.Errorw("DialContext", logx.Field("error", err.Error()))
		return
	}
	defer client.Close()

	chainID, err := client.ChainID(context.TODO())
	if err != nil {
		logx.Errorw("ChainID", logx.Field("error", err.Error()))
		return
	}
	w.chainId.Set(chainID)

	accountAddress := GetAddress(w.privateKey)

	nonce, err := client.PendingNonceAt(context.TODO(), accountAddress)
	if err != nil {
		logx.Errorw("PendingNonceAt", logx.Field("error", err.Error()))
		return
	}
	w.nonce = nonce
	// 需要启动一个新的协程，从而不会阻塞
	go func() {
		for {
			// 将worker注册到线程池
			workerPool.workerQueue <- w
			select {
			case job := <-w.jobQueue:
				//尝试5次
				tries := 5
			inner:
				for i := 0; i < tries; i++ {
					err := job.Do(w.privateKey, w.chainId, w.nonce)
					if err == nil {
						//成功发送加一
						w.nonce++
						break inner
					} else {
						//nonce太小
						if strings.Contains(err.Error(), NonceTooLow) {
							nonce, err := w.GetNonce()
							if err == nil {
								w.nonce = nonce
							} else {
								logx.Errorw("GetNonce", logx.Field("error", err.Error()))
							}
						}
						if strings.Contains(err.Error(), "txpool is full") {
							//如果发现交易池满，pending为0，
							nonce, err := w.GetNonce()
							if err == nil {
								logx.Errorw("GetNonce", logx.Field("new", nonce), logx.Field("old", w.nonce))
								w.nonce = nonce
							} else {
								logx.Errorw("GetNonce", logx.Field("error", err.Error()))
							}
							time.Sleep(time.Minute)
						}
						//交易池满等错误
						logx.Errorw("sendTransaction", logx.Field("error", err.Error()))
						time.Sleep(time.Second * 20)
					}
				}
			// 终止当前worker
			case <-w.stop:
				return
			}
		}
	}()
}
