package main

import (
	"fmt"
)

// ORDER 定义常量
const (
	ORDER = 3 // B+ 树的阶数
)

// BPlusTreeNode 定义 B+ 树的节点
type BPlusTreeNode struct {
	isLeaf   bool
	keys     []int
	children []*BPlusTreeNode
	next     *BPlusTreeNode
}

// BPlusTree 定义 B+ 树
type BPlusTree struct {
	root *BPlusTreeNode
}

// newBPlusTreeNode 创建一个新的 B+ 树节点
func newBPlusTreeNode(isLeaf bool) *BPlusTreeNode {
	return &BPlusTreeNode{
		isLeaf:   isLeaf,
		keys:     make([]int, 0),
		children: make([]*BPlusTreeNode, 0),
	}
}

// NewBPlusTree 创建一个新的 B+ 树
func NewBPlusTree() *BPlusTree {
	root := newBPlusTreeNode(true)
	return &BPlusTree{
		root: root,
	}
}

// Insert 插入键值到 B+ 树中
func (tree *BPlusTree) Insert(key int) {
	root := tree.root
	if len(root.keys) == 2*ORDER-1 {
		newRoot := newBPlusTreeNode(false)
		newRoot.children = append(newRoot.children, root)
		tree.splitChild(newRoot, 0)
		tree.root = newRoot
	}
	tree.insertNonFull(tree.root, key)
}

// splitChild 拆分子节点
func (tree *BPlusTree) splitChild(node *BPlusTreeNode, i int) {
	child := node.children[i]
	newNode := newBPlusTreeNode(child.isLeaf)
	node.keys = append(node.keys[:i], append([]int{child.keys[ORDER-1]}, node.keys[i:]...)...)
	node.children = append(node.children[:i+1], append([]*BPlusTreeNode{newNode}, node.children[i+1:]...)...)

	newNode.keys = append(newNode.keys, child.keys[ORDER:]...)
	child.keys = child.keys[:ORDER-1]

	if !child.isLeaf {
		newNode.children = append(newNode.children, child.children[ORDER:]...)
		child.children = child.children[:ORDER]
	} else {
		newNode.next = child.next
		child.next = newNode
	}
}

// insertNonFull 插入到非满节点
func (tree *BPlusTree) insertNonFull(node *BPlusTreeNode, key int) {
	i := len(node.keys) - 1
	if node.isLeaf {
		node.keys = append(node.keys, 0)
		for i >= 0 && key < node.keys[i] {
			node.keys[i+1] = node.keys[i]
			i--
		}
		node.keys[i+1] = key
	} else {
		for i >= 0 && key < node.keys[i] {
			i--
		}
		i++
		if len(node.children[i].keys) == 2*ORDER-1 {
			tree.splitChild(node, i)
			if key > node.keys[i] {
				i++
			}
		}
		tree.insertNonFull(node.children[i], key)
	}
}

// Search 查找键值在 B+ 树中的位置
func (tree *BPlusTree) Search(key int) *BPlusTreeNode {
	return tree.search(tree.root, key)
}

func (tree *BPlusTree) search(node *BPlusTreeNode, key int) *BPlusTreeNode {
	i := 0
	for i < len(node.keys) && key > node.keys[i] {
		i++
	}
	if node.isLeaf {
		if i < len(node.keys) && key == node.keys[i] {
			return node
		}
		return nil
	}
	if i < len(node.keys) && key == node.keys[i] {
		i++
	}
	return tree.search(node.children[i], key)
}

// Print 打印 B+ 树（调试用）
func (tree *BPlusTree) Print() {
	tree.print(tree.root, 0)
}

func (tree *BPlusTree) print(node *BPlusTreeNode, level int) {
	fmt.Printf("Level %d: ", level)
	for _, key := range node.keys {
		fmt.Printf("%d ", key)
	}
	fmt.Println()
	if !node.isLeaf {
		for _, child := range node.children {
			tree.print(child, level+1)
		}
	}
}

// 主函数
func main() {
	tree := NewBPlusTree()

	//keys := []int{10, 20, 5, 6, 12, 30, 7, 17}
	for key := range 10 {
		tree.Insert(key + 1)
	}

	tree.Print()

	searchKey := 6
	node := tree.Search(searchKey)
	if node != nil {
		fmt.Printf("Found key %d in node with keys: %v\n", searchKey, node.keys)
	} else {
		fmt.Printf("Key %d not found in the tree.\n", searchKey)
	}
}
