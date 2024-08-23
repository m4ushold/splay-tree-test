package splay

import (
	"fmt"
	"strings"
)

// Tree is weighted binary search tree which is based on Splay tree.
// original paper on Splay Trees: https://www.cs.cmu.edu/~sleator/papers/self-adjusting.pdf
type Tree[V Value] struct {
	root        *Node[V]
	rotateCount int
}

// NewTree creates a new instance of Tree.
func NewTree[V Value](root *Node[V]) *Tree[V] {
	return &Tree[V]{
		root: root,
	}
}

// Insert inserts the node at the last.
func (t *Tree[V]) Insert(node *Node[V]) *Node[V] {
	if t.root == nil {
		t.root = node
		return node
	}

	return t.InsertAfter(t.root, node)
}

// InsertAfter inserts the node after the given previous node.
func (t *Tree[V]) InsertAfter(prev *Node[V], node *Node[V]) *Node[V] {
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

// Splay moves the given node to the root.
func (t *Tree[V]) Splay(node *Node[V]) {
	if node == nil {
		return
	}

	for {
		if isLeftChild(node.parent) && isRightChild(node) {
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
			// zig
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

func (t *Tree[V]) updateTreeNode(node *Node[V]) {
	for node != nil {
		t.updateNode(node)
		node = node.parent
	}
}

func (t *Tree[V]) updateNode(node *Node[V]) {
	node.InitWeight()
	node.height = 1
	node.size = 1

	if node.left != nil {
		node.increaseWeight(node.leftWeight())
		node.size += node.left.size
		if node.left.height+1 > node.height {
			node.height = node.left.height + 1
		}
	}

	if node.right != nil {
		node.increaseWeight(node.rightWeight())
		node.size += node.right.size
		if node.right.height+1 > node.height {
			node.height = node.right.height + 1
		}
	}
}

// IndexOf Find the index of the given node.
func (t *Tree[V]) IndexOf(node *Node[V]) int {
	if node == nil || node != t.root && !node.hasLinks() {
		return -1
	}

	index := 0
	current := node
	var prev *Node[V]
	for current != nil {
		if prev == nil || prev == current.right {
			index += current.value.Len() + current.leftWeight()
		}
		prev = current
		current = current.parent
	}
	return index - node.value.Len()
}

func (t *Tree[V]) Kth(index int) *Node[V] {
	node := t.root
	for {
		for node.left != nil && node.left.size > index {
			node = node.left
		}
		if node.left != nil {
			index -= node.left.size
		}
		if index == 0 {
			break
		}
		index--
		node = node.right
	}
	t.Splay(node)
	return node
}

// Find returns the Node and offset of the given index.
func (t *Tree[V]) Find(index int) (*Node[V], int, error) {
	if t.root == nil {
		return nil, 0, nil
	}

	node := t.root
	offset := index
	for {
		if node.left != nil && offset <= node.leftWeight() {
			node = node.left
		} else if node.right != nil && node.leftWeight()+node.value.Len() < offset {
			offset -= node.leftWeight() + node.value.Len()
			node = node.right
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
func (t *Tree[V]) Delete(node *Node[V]) {
	t.Splay(node)

	leftTree := NewTree(node.left)
	if leftTree.root != nil {
		leftTree.root.parent = nil
	}

	rightTree := NewTree(node.right)
	if rightTree.root != nil {
		rightTree.root.parent = nil
	}

	if leftTree.root != nil {
		rightmost := leftTree.rightmost()
		leftTree.Splay(rightmost)
		leftTree.root.right = rightTree.root
		if rightTree.root != nil {
			rightTree.root.parent = leftTree.root
		}
		t.root = leftTree.root
	} else {
		t.root = rightTree.root
	}

	node.unlink()
	if t.root != nil {
		t.updateNode(t.root)
	}
}

// DeleteRange separates the range between given 2 boundaries from this Tree.
// This function separates the range to delete as a subtree
// by splaying outer boundary nodes.
// leftBoundary must exist because of 0-indexed initial dummy node of tree,
// but rightBoundary can be nil means range to delete includes the end of tree.
// Refer to the design document: ./design/range-deletion-in-slay-tree.md
func (t *Tree[V]) DeleteRange(leftBoundary, rightBoundary *Node[V]) {
	if rightBoundary == nil {
		t.Splay(leftBoundary)
		t.cutOffRight(leftBoundary)
		return
	}
	t.Splay(leftBoundary)
	t.Splay(rightBoundary)
	if rightBoundary.left != leftBoundary {
		t.rotateRight(leftBoundary)
	}
	t.cutOffRight(leftBoundary)
}

func (t *Tree[V]) cutOffRight(root *Node[V]) {
	traversePostorder(root.right, func(node *Node[V]) { node.InitWeight() })
	t.updateTreeNode(root)
}

func (t *Tree[V]) rotateLeft(pivot *Node[V]) {
	t.rotateCount++
	root := pivot.parent
	if root.parent != nil {
		if root == root.parent.left {
			root.parent.left = pivot
		} else {
			root.parent.right = pivot
		}
	} else {
		t.root = pivot
	}

	pivot.parent = root.parent

	root.right = pivot.left
	if root.right != nil {
		root.right.parent = root
	}

	pivot.left = root
	pivot.left.parent = pivot

	t.updateNode(root)
	t.updateNode(pivot)
}

func (t *Tree[V]) rotateRight(pivot *Node[V]) {
	t.rotateCount++
	root := pivot.parent
	if root.parent != nil {
		if root == root.parent.left {
			root.parent.left = pivot
		} else {
			root.parent.right = pivot
		}
	} else {
		t.root = pivot
	}
	pivot.parent = root.parent

	root.left = pivot.right
	if root.left != nil {
		root.left.parent = root
	}

	pivot.right = root
	pivot.right.parent = pivot

	t.updateNode(root)
	t.updateNode(pivot)
}

func (t *Tree[V]) rightmost() *Node[V] {
	node := t.root
	for node.right != nil {
		node = node.right
	}
	return node
}

// Len returns the size of this Tree.
func (t *Tree[V]) Len() int {
	if t.root == nil {
		return 0
	}

	return t.root.weight
}

func (t *Tree[V]) Height() int {
	if t.root == nil {
		return 0
	}
	return t.root.height
}

func (t *Tree[V]) RotateCount() int {
	if t.root == nil {
		return 0
	}
	return t.rotateCount
}

// debugging --------------------------------------------------------

// String returns a string containing node values.
func (t *Tree[V]) String() string {
	var builder strings.Builder
	traverseInOrder(t.root, func(node *Node[V]) {
		builder.WriteString(node.value.String())
	})
	return builder.String()
}

// ToTestString returns a string containing the metadata of the Node
// for debugging purpose.
func (t *Tree[V]) ToTestString() string {
	var builder strings.Builder

	traverseInOrder(t.root, func(node *Node[V]) {
		builder.WriteString(fmt.Sprintf(
			"[%d,%d,%d,%d]%s",
			node.weight,
			node.size,
			node.height,
			node.value.Len(),
			node.value.String(),
		))
	})
	return builder.String()
}
