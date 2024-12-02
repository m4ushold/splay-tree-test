package splay

type BlockingLinearOpTree[V Value] struct {
	Tree[V]
	linearCount int
	firstNode   *Node[V]
}

func NewBlockingLinearOpTree[V Value](root *Node[V]) *BlockingLinearOpTree[V] {
	return &BlockingLinearOpTree[V]{
		Tree:       *NewTree(root),
		linearCount: 0,
	}
}

// Insert inserts the node at the last.
func (t *BlockingLinearOpTree[V]) Insert(node *Node[V]) *Node[V] {
	if t.root == nil {
		t.root = node
		return node
	}

	return t.InsertAfter(t.root, node)
}

// InsertAfter inserts the node after the given previous node.
func (t *BlockingLinearOpTree[V]) InsertAfter(prev *Node[V], node *Node[V]) *Node[V] {
	if prev == t.root {
		t.linearCount++
		if t.linearCount == 1 {
			t.firstNode = node
		} else if t.linearCount > 500 {
			t.Splay(t.firstNode)
			t.linearCount = 0
		}
	} else {
		t.linearCount = 0
	}

	t.Splay(prev)
	t.root = node
	node.right = prev.right
	if prev.right != nil {
		prev.right.parent = node
	}
	node.left = prev
	prev.parent = node
	prev.right = nil

	t.updateNode(prev)
	t.updateNode(node)

	return node
}