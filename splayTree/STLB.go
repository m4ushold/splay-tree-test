package splay

import (
	"fmt"
	"math"
)

type STLB[V Value] struct {
	Tree[V]
	threshold int
}

func NewSTLB[V Value](root *Node[V], threshold int) *STLB[V] {
	return &STLB[V]{
		Tree:      *NewTree(root),
		threshold: threshold,
	}
}

func (t *STLB[V]) GetHeightDiff(node *Node[V]) int {
	return int(math.Abs(float64(node.leftHeight() - node.rightHeight())))
}

func (t *STLB[V]) Propagate(node *Node[V]) {
	if node == nil || !node.lazy {
		return
	}

	if node.left != nil && t.GetHeightDiff(node.left) > t.threshold {
		t.Balancing(node.left)
	}
	if node.right != nil && t.GetHeightDiff(node.right) > t.threshold {
		t.Balancing(node.right)
	}
}

func (t *STLB[V]) Balancing(node *Node[V]) {
	parent := node.parent
	node = t.KthElement(node, node.size/2)
	t.InternalSplay(node, parent)
}

func (t *STLB[V]) KthElement(node *Node[V], index int) *Node[V] {
	for {
		for node.left != nil && node.left.size > index {
			node = node.left
		}
		if node.left != nil {
			index -= node.left.size
		}
		if index == 0 {
			return node
		}
		index--
		node = node.right
	}
}

func (t *STLB[V]) Splay(node *Node[V]) {

	t.InternalSplay(node, nil)
}

// Splay moves the given node to the root.
func (t *STLB[V]) InternalSplay(node *Node[V], g *Node[V]) {
	if node == nil {
		return
	}

	for {
		if node.parent == g {
			if isLeftChild(node) {
				t.rotateRight(node)
			} else if isRightChild(node) {
				t.rotateLeft(node)
			}
			t.updateTreeNode(node)
			return
		} else if isLeftChild(node.parent) && isRightChild(node) {
			// zig-zag
			t.rotateLeft(node)
			t.rotateRight(node)
		} else if isRightChild(node.parent) && isLeftChild(node) {
			// zig-zag
			t.rotateRight(node)
			t.rotateLeft(node)
		} else if isLeftChild(node.parent) && isLeftChild(node) {
			// zig-zig
			t.rotateRight(node.parent)
			t.rotateRight(node)
		} else if isRightChild(node.parent) && isRightChild(node) {
			// zig-zig
			t.rotateLeft(node.parent)
			t.rotateLeft(node)
		} else {
			if isLeftChild(node) {
				t.rotateRight(node)
			} else if isRightChild(node) {
				t.rotateLeft(node)
			}
			t.updateTreeNode(node)
			return
		}
	}
}

// Find returns the Node and offset of the given index.
func (t *STLB[V]) Find(index int) (*Node[V], int, error) {
	if t.root == nil {
		return nil, 0, nil
	}

	node := t.root
	offset := index
	t.Propagate(node)
	for {
		if node.left != nil && offset <= node.leftWeight() {
			node = node.left
			t.Propagate(node)
		} else if node.right != nil && node.leftWeight()+node.value.Len() < offset {
			offset -= node.leftWeight() + node.value.Len()
			node = node.right
			t.Propagate(node)
		} else {
			offset -= node.leftWeight()
			break
		}
	}

	if offset > node.value.Len() {
		return nil, 0, fmt.Errorf("node length %d, index %d: %w", node.value.Len(), offset, ErrOutOfIndex)
	}

	t.Splay(node)
	return node, offset, nil
}

// Delete deletes the given node from this Tree.
func (t *STLB[V]) Delete(node *Node[V]) {
	if node == nil {
		return
	}

	t.Tree.Delete(node)
	if t.root != nil && int(math.Abs(float64(t.root.leftHeight()-t.root.rightHeight()))) > t.threshold {
		t.root.lazy = true
	}
}

func (t *STLB[V]) Insert(node *Node[V]) *Node[V] {
	node = t.Tree.Insert(node)
	if t.root != nil && int(math.Abs(float64(t.root.leftHeight()-t.root.rightHeight()))) > t.threshold {
		t.root.lazy = true
	}
	return node
}
