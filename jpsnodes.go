package main

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
	List []*JPSNode
}

func (this *JPSNodes) Len() int {
	return len(this.List)
}

func (this *JPSNodes) AddByPos(x, y int) {
	node := &JPSNode{}
	node.Pos = &Pos{x, y}
	this.List = append(this.List, node)
	return
}

func (this *JPSNodes) Add(node *JPSNode) {
	this.List = append(this.List, node)
	return
}

func (this *JPSNodes) RemoveByIndex(i int) int {
	if i >= this.Len() || i < 0 {
		return 0
	}
	this.List = append(this.List[:i], this.List[i+1:]...)
	return 1
}
