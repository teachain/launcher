package monitor

import (
	"github.com/zeromicro/go-zero/core/logx"
	"math/big"
	"os"
)

type BlockUpdater struct {
	filePath    string
	blockNumber string
}

func needToFix(localBlock string, remoteBlock string) bool {
	localBlockNumber := big.NewInt(0)

	remoteBlockNumber := big.NewInt(0)

	localBlockNumber.SetString(localBlock, 10)

	remoteBlockNumber.SetString(remoteBlock, 10)

	if localBlockNumber.Cmp(remoteBlockNumber) > 0 {
		return true
	}
	return false
}

func NewBlockUpdater(filename string, initBlockNumber string) (*BlockUpdater, error) {
	blockNumber := "0"
	if fileExists(filename) {
		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		if needToFix(string(data), initBlockNumber) {
			blockNumber = initBlockNumber
		} else {
			blockNumber = string(data)
		}
	} else {
		err := os.WriteFile(filename, []byte(initBlockNumber), os.ModePerm)
		if err != nil {
			return nil, err
		}
		blockNumber = initBlockNumber
	}
	b := &BlockUpdater{
		filePath:    filename,
		blockNumber: blockNumber,
	}
	return b, nil
}

func (b *BlockUpdater) OnBlock(number *big.Int) {
	logx.Debugw("OnBlock", logx.Field("number", number.String()))
	err := os.WriteFile(b.filePath, []byte(number.String()), os.ModePerm)
	if err != nil {
		logx.Errorw("WriteFile", logx.Field("error", err.Error()))
	}
}
func (b *BlockUpdater) GetBlockNumber() string {
	return b.blockNumber
}
