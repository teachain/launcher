package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/simplechain-org/client/ethclient"
	"github.com/teachain/launcher/internal/monitor"
	"github.com/teachain/launcher/pkg/onchain"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

var txCount *int = flag.Int("txCount", 1000, "-txCount=1000")
var dataSize *int = flag.Int("dataSize", 1024, "-dataSize=1024")
var httpEndpoint *string = flag.String("http", "http://192.168.4.31:18545", "-http=http://192.168.4.31:18545")
var wsEndpoint *string = flag.String("ws", "ws://192.168.4.31:18546", "-ws=ws://192.168.4.31:18546")
var noLimit *bool = flag.Bool("noLimit", false, "--noLimit=false")

func main() {
	flag.Parse()
	dir, err := os.Getwd()
	if err != nil {
		logx.Errorw("directory", logx.Field("error", err.Error()))
		return
	}
	filename := filepath.Join(dir, "blockNumber")

	chainClient, err := ethclient.DialContext(context.Background(), *httpEndpoint)
	if err != nil {
		logx.Errorw("chain node DialContext", logx.Field("error", err.Error()))
		return
	}
	blockNumber, err := chainClient.BlockNumber(context.Background())
	if err != nil {
		logx.Errorw("chain node BlockNumber", logx.Field("error", err.Error()))
		return
	}
	blockUpdater, err := monitor.NewBlockUpdater(filename, fmt.Sprintf("%d", blockNumber))
	if err != nil {
		logx.Errorw("NewBlockUpdater", logx.Field("error", err.Error()))
		return
	}
	config := onchain.ChainConfig{
		HttpEndpoint: *httpEndpoint,
		WsEndpoint:   *wsEndpoint,
		TxCount:      *txCount,
		DataSize:     *dataSize,
		LastBlock:    blockUpdater.GetBlockNumber(),
		NoLimit:      *noLimit,
	}
	manager := onchain.NewManager(context.Background(), config, blockUpdater)
	manager.Start()
	daemon()
}

func daemon() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Kill, os.Interrupt)
	<-ch
	logx.Infow("application exit", logx.Field("time", time.Now().Format(time.DateTime)))
}
