package onchain

import (
	"crypto/ecdsa"
	"github.com/zeromicro/go-zero/core/logx"
	"math/rand"
	"runtime"
	"time"
)

// PrepareAccount 准备上链账户
func (m *Manager) PrepareAccount() {
	//因为上链的worker个数就是runtime.NumCPU()
	accountCount := runtime.NumCPU()
	m.accounts = make([]*ecdsa.PrivateKey, 0)
	for i := 0; i < accountCount; i++ {
		privateKey, err := GenerateKey()
		if err != nil {
			i--
			m.Logger.Errorw("GenerateKey", logx.Field("error", err.Error()))
			continue
		}
		m.accounts = append(m.accounts, privateKey)
	}
}

func (m *Manager) randomAccount() *ecdsa.PrivateKey {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := random.Intn(len(m.accounts))
	return m.accounts[index]
}
