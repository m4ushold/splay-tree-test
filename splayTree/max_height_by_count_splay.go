package splay

type MaxHeightByCountSplayTree[V Value] struct {
	Tree[V]
	policy         func(int) int
	operationCount int
	splayCount     int
}

func NewMaxHeightByCountSplayTree[V Value](root *Node[V], policy func(int) int, splayCount int) *MaxHeightByCountSplayTree[V] {
	return &MaxHeightByCountSplayTree[V]{
		Tree:           *NewTree(root),
		policy:         policy,
		operationCount: 0,
		splayCount:     splayCount,
	}
}

func (t *MaxHeightByCountSplayTree[V]) MaxHeightSplay() {
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

func (t *MaxHeightByCountSplayTree[V]) countOperation() {
	t.operationCount++
	if t.root != nil {
		if t.operationCount >= t.policy(t.root.size) {
			for i := 0; i < t.splayCount; i++ {
				t.MaxHeightSplay()
			}
			t.operationCount = 0
		}
	}
}

func (t *MaxHeightByCountSplayTree[V]) Insert(node *Node[V]) *Node[V] {
	t.countOperation()
	return t.Tree.Insert(node)
}

func (t *MaxHeightByCountSplayTree[V]) Find(index int) (*Node[V], int, error) {
	t.countOperation()
	return t.Tree.Find(index)
}

func (t *MaxHeightByCountSplayTree[V]) Delete(node *Node[V]) {
	t.countOperation()
	t.Tree.Delete(node)
}
