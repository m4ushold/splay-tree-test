package splay

import (
	"math/rand"
)

type RandomBoundSplayTree[V Value] struct {
	Tree[V]
	upperBound int
	splayCount int
}

func NewRandomBoundSplayTree[V Value](root *Node[V], upperBound int, splayCount int) *RandomBoundSplayTree[V] {
	return &RandomBoundSplayTree[V]{
		Tree:       *NewTree(root),
		upperBound: upperBound,
		splayCount: splayCount,
	}
}

func (t *RandomBoundSplayTree[V]) checkHeight() {
	if t.root != nil {
		if t.root.height > t.upperBound {
			for i := 0; i < t.splayCount; i++ {
				t.Tree.Kth(rand.Intn(t.root.size))
			}
		}
	}
}

func (t *RandomBoundSplayTree[V]) Insert(node *Node[V]) *Node[V] {
	t.checkHeight()
	return t.Tree.Insert(node)
}

func (t *RandomBoundSplayTree[V]) Find(index int) (*Node[V], int, error) {
	t.checkHeight()
	return t.Tree.Find(index)
}

func (t *RandomBoundSplayTree[V]) Delete(node *Node[V]) {
	t.checkHeight()
	t.Tree.Delete(node)
}
