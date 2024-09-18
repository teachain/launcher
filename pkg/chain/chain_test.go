package chain

import (
	"context"
	"github.com/simplechain-org/client/common"
	"github.com/simplechain-org/client/ethclient"
	"testing"
)

func TestClient_TxPoolStatus(t *testing.T) {
	endpoint := "http://192.168.4.31:18545"
	client, err := NewClient(endpoint)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	poolStatus, err := client.TxPoolStatus()
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	pending := poolStatus.GetPending()
	t.Log(pending)
}

func TestClient_TransactionCountLatest(t *testing.T) {
	endpoint := "http://192.168.4.31:18545"
	client, err := NewClient(endpoint)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	account := common.HexToAddress("0x9ee44fea3e10a5c5a160eb1f61687eb464f5ac17")
	nonce, err := client.TransactionCountLatest(context.Background(), account)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	t.Log("nonce=", nonce)

	c, err := ethclient.DialContext(context.Background(), endpoint)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	n, err := c.PendingNonceAt(context.Background(), account)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	t.Log("pending nonce=", n)
}
