package onchain

import (
	"math/big"
)

// PriorityQueue 定义基于big.Int的优先队列
type PriorityQueue []*big.Int

// Len 实现heap.Interface的方法
func (pq *PriorityQueue) Len() int { return len(*pq) }

func (pq *PriorityQueue) Less(i, j int) bool {
	// 我们希望Pop给我们最小值而不是最大值，所以使用小于号
	return (*pq)[i].Cmp((*pq)[j]) < 0
}

func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*big.Int)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
