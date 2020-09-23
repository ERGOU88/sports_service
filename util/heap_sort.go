package util

import (
	"sort"
)


func PartialSort(data sort.Interface, m int) {
	// 建max-heap
	makeHeap(data, 0, m)
	minElemIdx := 0

	// 遍历后续的元素
	len := data.Len()
	for i := m; i < len; i++ {
		if data.Less(i, minElemIdx) {
			// 当这个后续元素比ResArr中最大的元素小，则替换
			data.Swap(i, minElemIdx)
			// 重新调整max-heap
			siftDown(data, minElemIdx, m, minElemIdx)
		}
	}

	// 对max-heap进行对排序
	heapsort(data, 0, m)
}

// 对data[lo, hi)内的元素进行siftDown处理，以满足heap序列的要求
func siftDown(data sort.Interface, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1 // 算出左子节点
		if child >= hi { // 判断是否超出范围
			break
		}

		// 找出值较大的子节点
		if child+1 < hi && data.Less(first+child, first+child+1) {
			child++
		}

		// 判断父节点是否小于最大的子节点，如果不小于则终止循环
		if !data.Less(first+root, first+child) {
			return
		}

		// 交换父节点和最大的子节点的值，继续对下一层做以上处理
		data.Swap(first+root, first+child)
		root = child
	}
}

// 将data[a, b)的数据排列为一个max-heap
func makeHeap(data sort.Interface, a, b int) {
	first := a
	hi := b - a

	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown(data, i, hi, first)
	}
}

// 将一个max-heap进行堆排序
func heapsort(data sort.Interface, a, b int) {
	first := a
	hi := b - a

	for i := hi - 1; i >= 0; i-- {
		data.Swap(first, first+i)
		siftDown(data, first, i, first)
	}
}
