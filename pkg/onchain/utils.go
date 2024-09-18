package onchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/simplechain-org/client/crypto"
)

const (
	NonceTooLow string = "nonce too low"
)

func GenerateKey() (*ecdsa.PrivateKey, error) {
	randBytes := make([]byte, 64)
	_, err := rand.Reader.Read(randBytes)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(randBytes)
	return ecdsa.GenerateKey(crypto.S256(), reader)
}

func Sha256(data []byte) string {
	// 创建一个新的hash.Hash对象，这里我们使用sha256.New()
	hashObj := sha256.New()

	// 写入要计算哈希的数据
	hashObj.Write(data)

	// Sum函数返回计算出的哈希值，它接受一个切片作为参数，
	// 该切片将被用作结果切片的开始部分（这里我们传入nil，表示不需要预分配的切片）
	// 返回的切片是哈希值的字节表示
	hashBytes := hashObj.Sum(nil)

	// 将字节切片转换为十六进制表示的字符串
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
