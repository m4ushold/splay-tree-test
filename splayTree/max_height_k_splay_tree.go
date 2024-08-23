package splay

import (
	"math"
)

type MaxHeightKSplayTree[V Value] struct {
	Tree[V]
	ratio      int
	splayCount int
}

func NewMaxHeightKSplayTree[V Value](root *Node[V], ratio int, splayCount int) *MaxHeightKSplayTree[V] {
	return &MaxHeightKSplayTree[V]{
		Tree:       *NewTree(root),
		ratio:      ratio,
		splayCount: splayCount,
	}
}

func (t *MaxHeightKSplayTree[V]) MaxHeightSplay() {
	node := t.Tree.root
	for node.height > 1 {
		if node.left != nil && node.left.height+1 == node.height {
			node = node.left
		} else {
			node = node.right
		}
	}
	t.Tree.Splay(node)
}

func (t *MaxHeightKSplayTree[V]) checkHeight() {
	if t.root != nil {
		if t.root.height > t.ratio*int(math.Log2(float64(t.root.size))) {
			for i := 0; i < t.splayCount; i++ {
				t.MaxHeightSplay()
			}
		}
	}
}

func (t *MaxHeightKSplayTree[V]) Insert(node *Node[V]) *Node[V] {
	t.checkHeight()
	return t.Tree.Insert(node)
}

func (t *MaxHeightKSplayTree[V]) Find(index int) (*Node[V], int, error) {
	t.checkHeight()
	return t.Tree.Find(index)
}

func (t *MaxHeightKSplayTree[V]) Delete(node *Node[V]) {
	t.checkHeight()
	t.Tree.Delete(node)
}
