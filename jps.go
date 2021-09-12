package main

import (
	"container/heap"
	"errors"
	"fmt"
)

type JPS struct {
	Map      [][]byte
	Start    *JPSNode
	End      *JPSNode
	Width    int
	Height   int
	OpenSet  *JPSNodes
	CloseSet *JPSNodes
	Path     []*Pos
}

func (this *JPS) Init(m [][]byte) {
	this.Height = len(m)
	if len(m) > 0 {
		this.Width = len(m[0])
	}
	this.Map = m
}

func (this *JPS) IsEnable(x, y int) bool {
	if x < 0 || y < 0 || x >= this.Width || y >= this.Height {
		return false
	}
	return this.Map[y][x] == 0
}

//此方法有切角问题
func (this *JPS) Neighbours(point *JPSNode) (list JPSNodes) {
	if point == nil {
		return list
	}
	parent := point.Parent
	if parent == nil {
		//起点
		for x := -1; x <= 1; x++ {
			for y := -1; y <= 1; y++ {
				if x == 0 && y == 0 {
					continue
				}
				if this.IsEnable(x+point.Pos.X, y+point.Pos.Y) {
					list.AddByPos(x+point.Pos.X, y+point.Pos.Y)
				}
			}
		}
		return list
	}

	xDirection := IBound(point.Pos.X-parent.Pos.X, -1, 1)
	yDirection := IBound(point.Pos.Y-parent.Pos.Y, -1, 1)
	if xDirection != 0 && yDirection != 0 {
		//对角方向
		neighbourForward := this.IsEnable(point.Pos.X, point.Pos.Y+yDirection)
		neighbourRight := this.IsEnable(point.Pos.X+xDirection, point.Pos.Y)
		neighbourLeft := this.IsEnable(point.Pos.X-xDirection, point.Pos.Y)
		neighbourBack := this.IsEnable(point.Pos.X, point.Pos.Y-yDirection)
		if neighbourForward {
			list.AddByPos(point.Pos.X, point.Pos.Y+yDirection)
		}
		if neighbourRight {
			list.AddByPos(point.Pos.X+xDirection, point.Pos.Y)
		}
		if (neighbourForward || neighbourRight) && this.IsEnable(point.Pos.X+xDirection, point.Pos.Y+yDirection) {
			list.AddByPos(point.Pos.X+xDirection, point.Pos.Y+yDirection)
		}
		//强迫邻居的处理,此算法对角移动会穿越阻挡
		if !neighbourLeft && neighbourForward {
			if this.IsEnable(point.Pos.X-xDirection, point.Pos.Y+yDirection) {
				list.AddByPos(point.Pos.X-xDirection, point.Pos.Y+yDirection)
			}
		}
		if !neighbourBack && neighbourRight {
			if this.IsEnable(point.Pos.X+xDirection, point.Pos.Y-yDirection) {
				list.AddByPos(point.Pos.X+xDirection, point.Pos.Y-yDirection)
			}
		}
	} else {
		if xDirection == 0 {
			//纵向
			if this.IsEnable(point.Pos.X, point.Pos.Y+yDirection) {
				list.AddByPos(point.Pos.X, point.Pos.Y+yDirection)
				//强迫邻居,此算法对角移动会穿越阻挡
				if !this.IsEnable(point.Pos.X+1, point.Pos.Y) && this.IsEnable(point.Pos.X+1, point.Pos.Y+yDirection) {
					list.AddByPos(point.Pos.X+1, point.Pos.Y+yDirection)
				}
				if !this.IsEnable(point.Pos.X-1, point.Pos.Y) && this.IsEnable(point.Pos.X-1, point.Pos.Y+yDirection) {
					list.AddByPos(point.Pos.X-1, point.Pos.Y+yDirection)
				}
			}
		} else {
			//横向
			if this.IsEnable(point.Pos.X+xDirection, point.Pos.Y) {
				list.AddByPos(point.Pos.X+xDirection, point.Pos.Y)
				//强迫邻居,此算法对角移动会穿越阻挡
				if !this.IsEnable(point.Pos.X, point.Pos.Y+1) && this.IsEnable(point.Pos.X+xDirection, point.Pos.Y+1) {
					list.AddByPos(point.Pos.X+xDirection, point.Pos.Y+1)
				}
				if !this.IsEnable(point.Pos.X, point.Pos.Y-1) && this.IsEnable(point.Pos.X+xDirection, point.Pos.Y-1) {
					list.AddByPos(point.Pos.X+xDirection, point.Pos.Y-1)
				}
			}
		}
	}
	return list
}

//此方法有切角问题
func (this *JPS) Jump(cur, pre *JPSNode, depth int) (node *JPSNode) {
	if !this.IsEnable(cur.Pos.X, cur.Pos.Y) {
		return nil
	}
	node = &JPSNode{}
	if depth == 0 || this.End.Pos.X == cur.Pos.X && this.End.Pos.Y == cur.Pos.Y { //递归到达指定深度，或到达目的地
		return &JPSNode{Pos: &Pos{cur.Pos.X, cur.Pos.Y}}
	}
	xDirection := IBound(cur.Pos.X-pre.Pos.X, -1, 1)
	yDirection := IBound(cur.Pos.Y-pre.Pos.Y, -1, 1)
	if xDirection != 0 && yDirection != 0 {
		if (this.IsEnable(cur.Pos.X+xDirection, cur.Pos.Y-yDirection) && !this.IsEnable(cur.Pos.X, cur.Pos.Y-yDirection)) ||
			(this.IsEnable(cur.Pos.X-xDirection, cur.Pos.Y+yDirection) && !this.IsEnable(cur.Pos.X-xDirection, cur.Pos.Y)) {
			/*
				| | | |
				| |x| |
				|p|1|j|

				|j| | |
				|1|x| |
				|p| | |
			*/
			return &JPSNode{Pos: &Pos{cur.Pos.X, cur.Pos.Y}}
		}
		//斜向移动的情况，横向或者纵向能找到跳点，则当前点为跳点
		//横向递归寻找强迫邻居
		if this.Jump(&JPSNode{Pos: &Pos{cur.Pos.X + xDirection, cur.Pos.Y}}, cur, depth-1) != nil {
			return &JPSNode{Pos: &Pos{cur.Pos.X, cur.Pos.Y}}
		}

		//纵向向递归寻找强迫邻居
		if this.Jump(&JPSNode{Pos: &Pos{cur.Pos.X, cur.Pos.Y + yDirection}}, cur, depth-1) != nil {
			return &JPSNode{Pos: &Pos{cur.Pos.X, cur.Pos.Y}}
		}
	} else if xDirection != 0 {
		//横向
		if (this.IsEnable(cur.Pos.X+xDirection, cur.Pos.Y+1) && !this.IsEnable(cur.Pos.X, cur.Pos.Y+1)) ||
			(this.IsEnable(cur.Pos.X+xDirection, cur.Pos.Y-1) && !this.IsEnable(cur.Pos.X, cur.Pos.Y-1)) {
			/*
				| |1|j|
				|p|x| |
				| | | |

				| | | |
				|p|x| |
				| |1|j|
			*/
			return &JPSNode{Pos: &Pos{cur.Pos.X, cur.Pos.Y}}
		}
	} else if yDirection != 0 {
		//纵向
		if (this.IsEnable(cur.Pos.X+1, cur.Pos.Y+yDirection) && !this.IsEnable(cur.Pos.X+1, cur.Pos.Y)) ||
			(this.IsEnable(cur.Pos.X-1, cur.Pos.Y+yDirection) && !this.IsEnable(cur.Pos.X-1, cur.Pos.Y)) {
			/*
				|j| | |
				|1|x| |
				| |p| |

				| | |j|
				| |x|1|
				| |p| |
			*/
			return &JPSNode{Pos: &Pos{cur.Pos.X, cur.Pos.Y}}
		}
	}
	//继续向当前方向前进
	return this.Jump(&JPSNode{Pos: &Pos{cur.Pos.X + xDirection, cur.Pos.Y + yDirection}}, cur, depth-1)
}

func (this *JPS) JPSNodeInCloseSet(jp *JPSNode) bool {
	for _, v := range this.CloseSet.Heap {
		if v.Pos.X == jp.Pos.X && v.Pos.Y == jp.Pos.Y {
			return true
		}
	}
	return false
}

func (this *JPS) JPSNodeInOpenSet(jp *JPSNode) *JPSNode {
	for _, v := range this.OpenSet.Heap {
		if v.Pos.X == jp.Pos.X && v.Pos.Y == jp.Pos.Y {
			return v
		}
	}
	return nil
}

func (this *JPS) CalcG(node1, node2 *JPSNode) int {
	if node1.Pos.X == node2.Pos.X {
		return IntAbs(node1.Pos.Y-node2.Pos.Y)*10 + node2.G
	} else if node1.Pos.Y == node2.Pos.Y {
		return IntAbs(node1.Pos.X-node2.Pos.X)*10 + node2.G
	} else {
		return IntAbs(node1.Pos.X-node2.Pos.X)*14 + node2.G
	}
}

func (this *JPS) CalcH(node1, node2 *JPSNode) int {
	return (IntAbs(node1.Pos.X-node2.Pos.X) + IntAbs(node1.Pos.Y-node2.Pos.Y)) * 10 //曼哈顿距离

}

//此方法有切角问题
func (this *JPS) ExtendRound(cur *JPSNode) {
	nbs := this.Neighbours(cur)
	for _, node := range nbs.Heap {
		jp := this.Jump(node, cur, 1000)
		if jp != nil {
			if this.JPSNodeInCloseSet(jp) {
				continue
			}
			jp.Parent = cur
			jp.G = this.CalcG(jp, cur)
			jp.H = this.CalcH(jp, cur)
			jp.F = jp.G + jp.H
			on := this.JPSNodeInOpenSet(jp)
			if on == nil {
				this.OpenSet.Add(jp)
				continue
			}
			if on.G > jp.G {
				on.Parent = cur
				on.G = jp.G
				on.F = jp.G + on.H
			}
		}
	}
	return
}

func (this *JPS) IsEnd(p *JPSNode) bool {
	if p == nil {
		return false
	}
	if p.Pos.X == this.End.Pos.X &&
		p.Pos.Y == this.End.Pos.Y {
		return true
	}
	return false
}

func (this *JPS) FindPath(startpos, endpos *Pos) error {
	this.Start, this.End = &JPSNode{Pos: startpos}, &JPSNode{Pos: endpos}
	this.Start.H = this.CalcH(this.Start, this.End)
	this.OpenSet = &JPSNodes{}
	this.CloseSet = &JPSNodes{}
	this.Path = nil
	this.OpenSet.Add(this.Start)
	for {
		if this.OpenSet.Len() == 0 {
			break
		}
		p, ok := heap.Pop(this.OpenSet).(*JPSNode)
		if !ok {
			break
		}
		if this.IsEnd(p) {
			this.MakePath(p)
			return nil
		}
		this.ExtendRound(p)
		this.CloseSet.Add(p)
	}
	return errors.New("FindPath: not find")
}

func (this *JPS) MakePath(node *JPSNode) {
	for node != nil {
		this.Path = append(this.Path, &Pos{node.Pos.X, node.Pos.Y})
		node = node.Parent
	}
	this.FixPath()
	return
}

func (this *JPS) FixPath() {
	//根据需要填充跳点中的路劲
	return
}

func (this *JPS) PrintPath() {
	fmt.Print("Path: ")
	for _, v := range this.Path {
		fmt.Print("[", v.X, ",", v.Y, "]")
	}
	fmt.Println()
	return
}
