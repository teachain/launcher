package onchain

import (
	"container/heap"
	"fmt"
	"math/big"
	"testing"
)

func TestPriorityQueue_Len(t *testing.T) {
	// 创建一个优先队列并给它一些值
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, big.NewInt(10))
	heap.Push(&pq, big.NewInt(60))
	heap.Push(&pq, big.NewInt(40))
	heap.Push(&pq, big.NewInt(20))
	heap.Push(&pq, big.NewInt(15))

	// 取出并打印所有值
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*big.Int)
		fmt.Println(item)
	}
}
