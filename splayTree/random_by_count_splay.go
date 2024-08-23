package splay

import (
	"math/rand"
)

type RandomByCountSplayTree[V Value] struct {
	Tree[V]
	policy         func(int) int
	operationCount int
	splayCount     int
}

func NewRandomByCountSplayTree[V Value](root *Node[V], policy func(int) int, splayCount int) *RandomByCountSplayTree[V] {
	return &RandomByCountSplayTree[V]{
		Tree:           *NewTree(root),
		policy:         policy,
		operationCount: 0,
		splayCount:     splayCount,
	}
}

func (t *RandomByCountSplayTree[V]) countOperation() {
	t.operationCount++
	if t.root != nil {
		if t.operationCount >= t.policy(t.root.size) {
			for i := 0; i < t.splayCount; i++ {
				t.Tree.Kth(rand.Intn(t.root.size))
			}
			t.operationCount = 0
		}
	}
}

func (t *RandomByCountSplayTree[V]) Insert(node *Node[V]) *Node[V] {
	t.countOperation()
	return t.Tree.Insert(node)
}

func (t *RandomByCountSplayTree[V]) Find(index int) (*Node[V], int, error) {
	t.countOperation()
	return t.Tree.Find(index)
}

func (t *RandomByCountSplayTree[V]) Delete(node *Node[V]) {
	t.countOperation()
	t.Tree.Delete(node)
}
