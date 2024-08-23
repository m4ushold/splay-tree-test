package splay

type BasicSplayTree[V Value] struct {
	Tree[V]
}

func NewBasicSplayTree[V Value](root *Node[V]) *BasicSplayTree[V] {
	return &BasicSplayTree[V]{
		Tree: *NewTree(root),
	}
}
