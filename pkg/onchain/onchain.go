package onchain

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"github.com/google/uuid"
	"github.com/simplechain-org/client/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
	"math/big"
	"time"
)

// StartConsumer 上链模块
func (m *Manager) StartConsumer(txCount int, dataSize int) {
	if !m.chainConfig.NoLimit {
		for i := 0; i < txCount; i++ {
			data, err := Generate(dataSize)
			if err != nil {
				i--
				fmt.Println(err.Error())
			}
			m.batchOnChain.Enqueue(&TaskOnChain{
				Data:         data,
				MsgId:        uuid.New().String(),
				To:           m.randomAccount(),
				HttpEndpoint: m.chainConfig.HttpEndpoint,
				txCh:         m.txCh,
			})
		}
		m.startedAt = time.Now()
	} else {
		go func() {
			m.startedAt = time.Now()
			for {
				data, err := Generate(dataSize)
				if err != nil {
					fmt.Println(err.Error())
					time.Sleep(time.Minute)
					continue
				}
				m.batchOnChain.Enqueue(&TaskOnChain{
					Data:         data,
					MsgId:        uuid.New().String(),
					To:           m.randomAccount(),
					HttpEndpoint: m.chainConfig.HttpEndpoint,
					txCh:         m.txCh,
				})
			}
		}()
	}
}
func Generate(dataSize int) ([]byte, error) {
	data := make([]byte, dataSize)
	n, err := rand.Read(data)
	if err != nil {
		return nil, err
	}
	return data[:n], nil
}

type TaskOnChain struct {
	Data         []byte
	MsgId        string
	To           *ecdsa.PrivateKey
	HttpEndpoint string
	txCh         chan string
}

func (t *TaskOnChain) SendTransaction(privateKey *ecdsa.PrivateKey, chainID *big.Int, nonce uint64, data []byte) (string, error) {
	client, err := ethclient.DialContext(context.Background(), t.HttpEndpoint)
	if err != nil {
		return "", err
	}
	defer client.Close()
	accountAddress := GetAddress(t.To)
	signTx, err := BuildTx(client, privateKey, chainID, nonce, accountAddress, data)
	if err != nil {
		return "", err
	}
	//上链
	err = client.SendTransaction(context.Background(), signTx)
	if err != nil {
		return "", err
	}
	return signTx.Hash().String(), nil
}

func (t *TaskOnChain) Do(privateKey *ecdsa.PrivateKey, chainID *big.Int, nonce uint64) error {
	//上链操作（发送交易）
	txHash, err := t.SendTransaction(privateKey, chainID, nonce, t.Data)
	if err != nil {
		logx.Errorw("SendTransaction", logx.Field("error", err.Error()))
		//特别要注意这里，别把err给隐藏了，一定要返回给上层进行error的处理，涉及到nonce的连续性
		return err
	} else {
		t.txCh <- txHash
		return nil
	}
}
