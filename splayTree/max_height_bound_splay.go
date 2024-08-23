package splay

type MaxHeightBoundSplayTree[V Value] struct {
	Tree[V]
	upperBound int
	splayCount int
}

func NewMaxHeightBoundSplayTree[V Value](root *Node[V], upperBound int, splayCount int) *MaxHeightBoundSplayTree[V] {
	return &MaxHeightBoundSplayTree[V]{
		Tree:       *NewTree(root),
		upperBound: upperBound,
		splayCount: splayCount,
	}
}

func (t *MaxHeightBoundSplayTree[V]) MaxHeightSplay() {
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

func (t *MaxHeightBoundSplayTree[V]) checkHeight() {
	if t.root != nil {
		if t.root.height > t.upperBound {
			for i := 0; i < t.splayCount; i++ {
				t.MaxHeightSplay()
			}
		}
	}
}

func (t *MaxHeightBoundSplayTree[V]) Insert(node *Node[V]) *Node[V] {
	t.checkHeight()
	return t.Tree.Insert(node)
}

func (t *MaxHeightBoundSplayTree[V]) Find(index int) (*Node[V], int, error) {
	t.checkHeight()
	return t.Tree.Find(index)
}

func (t *MaxHeightBoundSplayTree[V]) Delete(node *Node[V]) {
	t.checkHeight()
	t.Tree.Delete(node)
}
