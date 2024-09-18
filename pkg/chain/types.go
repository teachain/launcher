package chain

import (
	"strconv"
)

type TxPoolStatusResponse struct {
	Pending string `json:"pending"`
	Queued  string `json:"queued"`
}

func (t *TxPoolStatusResponse) GetPending() int64 {
	pending := t.Pending
	if has0xPrefix(pending) {
		pending = pending[2:]
	}
	result, err := strconv.ParseInt(pending, 16, 64)
	if err != nil {
		return 0
	}
	return result
}
func (t *TxPoolStatusResponse) GetQueued() int64 {
	queued := t.Queued
	if has0xPrefix(queued) {
		queued = queued[2:]
	}
	result, err := strconv.ParseInt(queued, 16, 64)
	if err != nil {
		return 0
	}
	return result
}
