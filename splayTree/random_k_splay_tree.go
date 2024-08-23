package splay

import (
	"math"
	"math/rand"
)

type RandomKSplayTree[V Value] struct {
	Tree[V]
	ratio      int
	splayCount int
}

func NewRandomKSplayTree[V Value](root *Node[V], ratio int, splayCount int) *RandomKSplayTree[V] {
	return &RandomKSplayTree[V]{
		Tree:       *NewTree(root),
		ratio:      ratio,
		splayCount: splayCount,
	}
}

func (t *RandomKSplayTree[V]) checkHeight() {
	if t.root != nil {
		if t.root.height > t.ratio*int(math.Log2(float64(t.root.size))) {
			for i := 0; i < t.splayCount; i++ {
				t.Kth(rand.Intn(t.root.size))
			}
		}
	}
}

func (t *RandomKSplayTree[V]) Insert(node *Node[V]) *Node[V] {
	t.checkHeight()
	return t.Tree.Insert(node)
}

func (t *RandomKSplayTree[V]) Find(index int) (*Node[V], int, error) {
	t.checkHeight()
	return t.Tree.Find(index)
}

func (t *RandomKSplayTree[V]) Delete(node *Node[V]) {
	t.checkHeight()
	t.Tree.Delete(node)
}
