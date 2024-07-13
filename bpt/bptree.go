package tree

import (
	"fmt"
)

// BPTree 存储节点
type BPTree struct {
	Root      *BPTreeNode //根节点
	size      int         //统计规模
	m         int         //度
	splitShow bool        //节点分列是否显示 默认不显示
	headNode  *BPTreeNode
	curNode   *BPTreeNode
}

// BPTreeNode 存储node 节点
type BPTreeNode struct {
	Parent   *BPTreeNode   //父节点
	Entries  []*Item       //存储实体
	Children []*BPTreeNode //叶子节点 列表
	Leaf     bool          //标记是不是叶子节点
	Next     *BPTreeNode
}

// Item 存储的实体数据
type Item struct {
	Key     int
	Val     interface{}
	deleted bool
}

func (i *Item) info() string {
	return fmt.Sprintf("(%d:deleted:%v)", i.Key, i.deleted)
}

// NewBPTree 构建新树
func NewBPTree(order int) *BPTree {
	if order < 3 {
		order = 3
	}
	return &BPTree{m: order}
}

// Put 进行put操作 存储数据
func (t *BPTree) Put(entry *Item) {
	// 如果根节点就是空的 直接结束
	if t.Root == nil {
		t.Root = &BPTreeNode{Entries: []*Item{entry}, Leaf: true}
		t.curNode = t.Root
		t.headNode = t.Root
		t.size++
		return
	}
	// 根节点不为空进行插入操作 插入成功 更新 size
	if t.insert(t.Root, entry) {
		t.size++
	}
}

// 进行插入操作
func (t *BPTree) insert(node *BPTreeNode, entry *Item) bool {
	// 如果是叶子节点 那就进行叶子节点插入工作
	// fmt.Println("insert ndoe:", node, entry)
	if t.isLeaf(node) {
		return t.insertIntoLeaf(node, entry)
	}
	// 如果不是叶子节点 那就继续向下查找
	// 因为b树的叶子插入都是在叶子节点进行的
	return t.insertIntoInternal(node, entry)
}

func setNewEntry(entires []*Item, index int, entry *Item) {
	entires[index].Key = entry.Key
	entires[index].Val = entry.Val
	entires[index].deleted = false
}

// 叶子节点的插入
// 找到对应的叶子节点 吧
// 返回如果是false 表示是修改 如果是true表示插入
func (t *BPTree) insertIntoLeaf(node *BPTreeNode, entry *Item) bool {
	// 进行节点查找
	idx, ok := t.search(node, entry.Key)
	if ok {
		// node.Entries[idx] = entry
		// return false
		setNewEntry(node.Entries, idx, entry)
		return false
	}
	// 开辟新的空间 多一个nil
	node.Entries = append(node.Entries, nil)
	// copy数据
	copy(node.Entries[idx+1:], node.Entries[idx:])
	// 在对应位置直接插入数据即可
	node.Entries[idx] = entry
	// 对于当前node节点 进行split操作试探

	t.split(node)
	return true
}

// 进行split操作
func (t *BPTree) split(node *BPTreeNode) {
	// 检查是否达到split操作要求
	if !t.shouldSplit(node) {
		return
	}
	if t.splitShow {
		// 如果达到split标准的话
		fmt.Println("node:", node, "start split")
	}
	// 如果是根节点的话 进行特殊处理
	if node == t.Root {
		t.splitRoot()
		return
	}
	// 否则进行正常处理即可
	t.splitNonRoot(node)
}

// 分割根节点
func (t *BPTree) splitRoot() {
	// 找到mid索引
	mid := t.middle()
	// 左边节点
	left := &BPTreeNode{
		Entries: append([]*Item(nil), t.Root.Entries[:mid]...),
		Leaf:    true,
	}

	right := &BPTreeNode{
		Entries: append([]*Item(nil), t.Root.Entries[mid:]...),
		Leaf:    true,
	}
	if t.headNode == t.Root {
		t.headNode = left
	}
	left.Next = right
	t.curNode = right
	// 判断根节点是否叶子节点
	// 如果根节点不是叶子节点的话
	if !t.isLeaf(t.Root) {
		left.Children = append([]*BPTreeNode(nil), t.Root.Children[:mid+1]...)
		right.Children = append([]*BPTreeNode(nil), t.Root.Children[mid+1:]...)
		// 分别为左右节点 设置父亲节点
		setParent(left.Children, left)
		setParent(right.Children, right)
	}
	// 产生新的根节点
	newRoot := &BPTreeNode{
		Entries:  []*Item{t.Root.Entries[mid]},
		Children: []*BPTreeNode{left, right},
	}
	// 更新根节点
	left.Parent = newRoot
	right.Parent = newRoot
	t.Root = newRoot
}

// 分割非根节点
func (t *BPTree) splitNonRoot(node *BPTreeNode) {
	// 找到mid索引下标
	middle := t.middle()
	// 找到父节点 因为对于当前节点需要进行拆分了
	// 是从最下面开始拆分的
	parent := node.Parent
	parent.Leaf = false
	// 找到左边部分
	left := &BPTreeNode{Entries: append([]*Item(nil), node.Entries[:middle]...), Parent: parent}
	// 找到右边部分
	// right := &BPTreeNode{Entries: append([]*Item(nil), node.Entries[middle+1:]...), Parent: parent}
	fmt.Println("leaf", t.isLeaf(node), node, node.Children, len(node.Entries))
	fmt.Println("root", t.Root)
	var right *BPTreeNode
	if t.isLeaf(node) {
		right = &BPTreeNode{Entries: append([]*Item(nil), node.Entries[middle:]...), Parent: parent}
	} else {
		right = &BPTreeNode{Entries: append([]*Item(nil), node.Entries[middle+1:]...), Parent: parent}
	}
	if node == t.curNode {
		t.curNode = left
		left.Next = right
	}
	t.curNode = right
	// 如果当前node节点 不是叶子节点的话
	if !t.isLeaf(node) {
		// 对于左边的孩子处理
		left.Children = append([]*BPTreeNode(nil), node.Children[:middle+1]...)
		// 对于右边孩子进行处理
		right.Children = append([]*BPTreeNode(nil), node.Children[middle+1:]...)
		// 分别设置父节点
		setParent(left.Children, left)
		setParent(right.Children, right)
	}

	// 找到需要插入的位置
	insertPosition, _ := t.search(parent, node.Entries[middle].Key)
	// 进行位置补充 需要把中间的 这个提上去
	parent.Entries = append(parent.Entries, nil)

	copy(parent.Entries[insertPosition+1:], parent.Entries[insertPosition:])
	// 将中间的这个提上去
	parent.Entries[insertPosition] = node.Entries[middle]

	// 左边孩子节点
	parent.Children[insertPosition] = left

	parent.Children = append(parent.Children, nil)
	// 为什么是+2呢 因为之前已经进行了+1 现在要在这个后面继续补一个
	// 需要移动一下 新数据已经填入了
	copy(parent.Children[insertPosition+2:], parent.Children[insertPosition+1:])
	// 右边孩子节点
	parent.Children[insertPosition+1] = right

	fmt.Println("root", t.Root)
	parent.Leaf = false
	// 继续往上走 开始递归上传即可
	t.split(parent)

}

// 设置父亲节点
func setParent(nodes []*BPTreeNode, parent *BPTreeNode) {
	for _, node := range nodes {
		node.Parent = parent
	}
}

// 节点个数大于最大个数的话
func (t *BPTree) shouldSplit(node *BPTreeNode) bool {
	return len(node.Entries) > t.maxEntries()
}
func (t *BPTree) Empty() bool {
	return t.size == 0
}
func (t *BPTree) middle() int {
	return (t.m - 1) / 2
}
func (t *BPTree) maxChildren() int {
	return t.m
}
func (t *BPTree) minChildren() int {
	return (t.m + 1) / 2
}
func (t *BPTree) maxEntries() int {
	return t.maxChildren() - 1
}
func (t *BPTree) minEntries() int {
	return t.minChildren() - 1
}
func (t *BPTree) isLeaf(node *BPTreeNode) bool {
	return node.Leaf //|| (len(node.Children) > 0 && t.isLeaf(node.Children[0]))
	// return len(node.Children) == 0
}

// 不断地进行二分查找 进行搜索 知道找到最终的叶子节点所在位置
// 节点一定是
func (t *BPTree) insertIntoInternal(node *BPTreeNode, entry *Item) bool {
	idx, ok := t.search(node, entry.Key)
	if ok {
		setNewEntry(node.Entries, idx, entry)
		return false
	}
	// todo
	// if idx >= len(node.Children) {
	// 	idx = min(idx, len(node.Children)-1)
	// }
	fmt.Println("node", node.Children, node.Entries, idx, ok)

	return t.insert(node.Children[idx], entry)
}

// 查找操作 简单的二分 靠左查询
func (t *BPTree) search(node *BPTreeNode, key int) (int, bool) {
	left := 0
	right := len(node.Entries) - 1
	ans := -1

	for left <= right {
		mid := (left + right) / 2
		if node.Entries[mid].Key == key {
			return mid, true
		} else if node.Entries[mid].Key < key {
			ans = mid
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return ans + 1, false
}

// 递归寻找
func (t *BPTree) searchRecur(node *BPTreeNode, key int) (st *BPTreeNode, idx int, ok bool) {
	if t.Empty() {
		return nil, -1, false
	}
	st = node
	for {
		idx, ok = t.search(node, key)
		if ok {
			return node, idx, true
		}
		if t.isLeaf(node) {
			return nil, -1, false
		}
		node = node.Children[idx]
	}
}

// Get 实现get操作
func (t *BPTree) Get(entry *Item) bool {
	node, idx, ok := t.searchRecur(t.Root, entry.Key)
	entry.Val = node.Entries[idx]
	return ok
}

// Remove 进行移除操作
func (t *BPTree) Remove(entry *Item) bool {
	// 查找数据是否存在
	node, idx, ok := t.searchRecur(t.Root, entry.Key)
	// 如果存在节点的话 那就需要对于节点进行删除操作了
	// 并且更新数据规模
	if ok {
		// 对于数据进行删除
		t.delete(node, idx)
		t.size--
	}
	return ok
}

// 实现真正意义上的删除操作
func (t *BPTree) delete(node *BPTreeNode, idx int) {
	// 如果是叶子节点 删除叶子节点
	// 否则删除内部节点
	node.Entries[idx].deleted = true
}

// Inorder 中序遍历
func (t *BPTree) Inorder() []*Item {
	if t.Root == nil {
		return []*Item{}
	}
	return inorder(t.Root)
}

func inorder(node *BPTreeNode) []*Item {
	if node == nil {
		return []*Item{}
	}
	result := []*Item{}
	// 遍历节点中的每个Item和子节点
	for i := 0; i < len(node.Entries); i++ {
		// 先递归遍历当前Item之前的子节点
		if i < len(node.Children) {
			result = append(result, inorder(node.Children[i])...)
		}
		// 访问当前的Item
		result = append(result, node.Entries[i])
	}
	// 最后递归遍历最后一个子节点
	if len(node.Children) > len(node.Entries) {
		result = append(result, inorder(node.Children[len(node.Entries)])...)
	}
	return result
}

// LevelOrderTraversal 层次遍历函数
func (t *BPTree) LevelOrderTraversal() (ans [][]*Item) {
	return levelOrderTraversal(t.Root)
}
func levelOrderTraversal(root *BPTreeNode) (ans [][]*Item) {
	if root == nil {
		return
	}
	queue := []*BPTreeNode{root}
	for i := 0; i < len(queue); i++ {
		ans = append(ans, []*Item{})
		var pueue []*BPTreeNode
		for j := 0; j < len(queue); j++ {
			node := queue[j]
			ans[i] = append(ans[i], node.Entries...)
			// 将当前节点的子节点加入队列
			for _, child := range node.Children {
				pueue = append(pueue, child)
			}
		}
		queue = pueue
	}
	return
}
