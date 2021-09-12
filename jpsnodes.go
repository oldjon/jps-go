package main

import "container/heap"

type Pos struct {
	X int
	Y int
}

type JPSNode struct {
	Pos    *Pos
	Parent *JPSNode
	G      int
	H      int
	F      int
}

type JPSNodes struct {
	Heap []*JPSNode
}

func (this *JPSNodes) Len() int {
	return len(this.Heap)
}

func (this *JPSNodes) AddByPos(x, y int) {
	node := &JPSNode{}
	node.Pos = &Pos{x, y}
	heap.Push(this, node)
	return
}

func (this *JPSNodes) Add(node *JPSNode) {
	heap.Push(this, node)
	return
}

func (this *JPSNodes) RemoveByIndex(i int) int {
	if i >= this.Len() || i < 0 {
		return 0
	}
	this.Heap = append(this.Heap[:i], this.Heap[i+1:]...)
	return 1
}

func (this *JPSNodes) Less(i, j int) bool {
	return this.Heap[i].F < this.Heap[j].F
}

func (this *JPSNodes) Swap(i, j int) {
	this.Heap[i], this.Heap[j] = this.Heap[j], this.Heap[i]
	return
}

func (this *JPSNodes) Push(x interface{}) {
	this.Heap = append(this.Heap, x.(*JPSNode))
	return
}

func (this *JPSNodes) Pop() interface{} {
	n := this.Len()
	x := this.Heap[n-1]
	this.Heap = this.Heap[0 : n-1]
	return x
}
