package splay

import "fmt"

// ErrOutOfIndex is returned when the given index is out of index.
var ErrOutOfIndex = fmt.Errorf("out of index")

// Value represents the data stored in the nodes of Tree.
type Value interface {
	Len() int
	String() string
}

// Node is a node of Tree.
type Node[V Value] struct {
	value  V
	weight int
	height int
	size   int
	lazy   bool

	left   *Node[V]
	right  *Node[V]
	parent *Node[V]
}

// NewNode creates a new instance of Node.
func NewNode[V Value](value V) *Node[V] {
	n := &Node[V]{
		value:  value,
		height: 1,
		size:   1,
		lazy:   false,
	}
	n.InitWeight()
	return n
}

// Value returns the value of this Node.
func (t *Node[V]) Value() V {
	return t.value
}

func (t *Node[V]) leftWeight() int {
	if t.left == nil {
		return 0
	}
	return t.left.weight
}

func (t *Node[V]) rightWeight() int {
	if t.right == nil {
		return 0
	}
	return t.right.weight
}

func (t *Node[V]) leftHeight() int {
	if t.left == nil {
		return 0
	}
	return t.left.height
}

func (t *Node[V]) rightHeight() int {
	if t.right == nil {
		return 0
	}
	return t.right.height
}

// InitWeight sets initial weight of this node.
func (t *Node[V]) InitWeight() {
	t.weight = t.value.Len()
}

func (t *Node[V]) increaseWeight(weight int) {
	t.weight += weight
}

func (t *Node[V]) unlink() {
	t.parent = nil
	t.right = nil
	t.left = nil
}

func (t *Node[V]) hasLinks() bool {
	return t.parent != nil || t.left != nil || t.right != nil
}

func isLeftChild[V Value](node *Node[V]) bool {
	return node != nil && node.parent != nil && node.parent.left == node
}

func isRightChild[V Value](node *Node[V]) bool {
	return node != nil && node.parent != nil && node.parent.right == node
}

func traverseInOrder[V Value](node *Node[V], callback func(node *Node[V])) {
	if node == nil {
		return
	}

	traverseInOrder(node.left, callback)
	callback(node)
	traverseInOrder(node.right, callback)
}

func traversePostorder[V Value](node *Node[V], callback func(node *Node[V])) {
	if node == nil {
		return
	}

	traversePostorder(node.left, callback)
	traversePostorder(node.right, callback)
	callback(node)
}
