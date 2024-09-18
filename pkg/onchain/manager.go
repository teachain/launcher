package onchain

import (
	"context"
	"crypto/ecdsa"
	"github.com/teachain/treasure/accelerate"
	"github.com/zeromicro/go-zero/core/logx"
	"math/big"
	"runtime"
	"time"
)

type BlockUpdater interface {
	OnBlock(number *big.Int)
}

type Manager struct {
	accounts []*ecdsa.PrivateKey
	logx.Logger
	ctx           context.Context
	dispatcher    *accelerate.Dispatcher
	handleNumbers chan *big.Int
	batchOnChain  *Dispatcher
	chainConfig   ChainConfig
	blockUpdater  BlockUpdater
	txCh          chan string
	onChainCount  int
	onChainTxCh   chan string
	startedAt     time.Time
}

func NewManager(ctx context.Context, chainConfig ChainConfig, blockUpdater BlockUpdater) *Manager {
	return &Manager{
		Logger:        logx.WithContext(ctx),
		dispatcher:    accelerate.NewDispatcher(runtime.NumCPU(), chainConfig.TxCount),
		handleNumbers: make(chan *big.Int, chainConfig.TxCount),
		batchOnChain:  NewDispatcher(runtime.NumCPU(), chainConfig.TxCount),
		chainConfig:   chainConfig,
		blockUpdater:  blockUpdater,
		txCh:          make(chan string, chainConfig.TxCount),
		onChainTxCh:   make(chan string, chainConfig.TxCount),
		startedAt:     time.Now(),
	}
}

func (m *Manager) Start() {
	m.dispatcher.Start()
	//准备上链的账户
	m.PrepareAccount()
	m.TransactionSent()
	//从消息队列里取数据，然后调用上链接口进行上链操作
	m.StartConsumer(m.chainConfig.TxCount, m.chainConfig.DataSize)
	//更新上链状态
	m.UpdateOnChainStatus(m.chainConfig.LastBlock)
	m.batchOnChain.Start(m.chainConfig.HttpEndpoint, m.accounts)
}

type ChainConfig struct {
	HttpEndpoint string `json:"httpEndpoint" yaml:"httpEndpoint"`
	WsEndpoint   string `json:"wsEndpoint" yaml:"wsEndpoint"`
	TxCount      int
	DataSize     int
	LastBlock    string
	NoLimit      bool
}
