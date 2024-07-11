package tree

const (
	ORDER = 3 //阶数
)

// BtreeNode 定义 B+ 树的节点
type Node struct {
	isLeaf   bool
	keys     []int
	children []*Node
	next     *Node
}

// Tree 定义 B+ 树
type Tree struct {
	root  *Node
	order int
}

// 创建一个新的 B+ 树节点
func newNode(isLeaf bool) *Node {
	return &Node{
		isLeaf:   isLeaf,
		keys:     make([]int, 0),
		children: make([]*Node, 0),
	}
}

// NewBTree 创建一个新的 B+ 树
func NewTree(order int) *Tree {
	root := newNode(true)
	return &Tree{
		root:  root,
		order: order,
	}
}
func (t *Tree) Insert(key int, value interface{}) {
	root := t.root
	if len(root.keys) < t.getMaxDegree() {

	}
}
func (t *Tree) insertNotFull(node *Node, key int) {
}

// 一次存储的最大节点个数
func (t *Tree) getMaxDegree() int {
	return t.order - 1
}
func (t *Tree) getMidIndex() int {
	return t.order / 2
}
