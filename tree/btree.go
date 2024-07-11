package tree

import (
	"fmt"
)

// BTree 存储节点
type BTree struct {
	Root      *BTreeNode //根节点
	size      int        //统计规模
	m         int        //度
	splitShow bool       //节点分列是否显示 默认不显示
}

// BTreeNode 存储node 节点
type BTreeNode struct {
	Parent   *BTreeNode   //父节点
	Entries  []*Item      //存储实体
	Children []*BTreeNode //叶子节点 列表
}

// Item 存储的实体数据
type Item struct {
	Key int
	Val interface{}
}

func (i *Item) info() string {
	return fmt.Sprintf("(%d:%v)", i.Key, i.Val)
}

// NewBTree 构建新树
func NewBTree(order int) *BTree {
	if order < 3 {
		order = 3
	}
	return &BTree{m: order}
}

// Put 进行put操作 存储数据
func (t *BTree) Put(entry *Item) {
	// 如果根节点就是空的 直接结束
	if t.Root == nil {
		t.Root = &BTreeNode{Entries: []*Item{entry}}
		t.size++
		return
	}
	// 根节点不为空进行插入操作 插入成功 更新 size
	if t.insert(t.Root, entry) {
		t.size++
	}
}

// 进行插入操作
func (t *BTree) insert(node *BTreeNode, entry *Item) bool {
	// 如果是叶子节点 那就进行叶子节点插入工作
	if t.isLeaf(node) {
		return t.insertIntoLeaf(node, entry)
	}
	// 如果不是叶子节点 那就继续向下查找
	// 因为b树的叶子插入都是在叶子节点进行的
	return t.insertIntoInternal(node, entry)
}

// 叶子节点的插入
// 找到对应的叶子节点 吧
func (t *BTree) insertIntoLeaf(node *BTreeNode, entry *Item) bool {
	// 进行节点查找
	idx, ok := t.search(node, entry.Key)
	if ok {
		node.Entries[idx] = entry
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
func (t *BTree) split(node *BTreeNode) {
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
func (t *BTree) splitRoot() {
	// 找到mid索引
	mid := t.middle()
	// 左边节点
	left := &BTreeNode{
		Entries: append([]*Item(nil), t.Root.Entries[:mid]...),
	}
	// 右边节点
	right := &BTreeNode{
		Entries: append([]*Item(nil), t.Root.Entries[mid+1:]...),
	}
	// 判断根节点是否叶子节点
	// 如果根节点不是叶子节点的话
	if !t.isLeaf(t.Root) {
		left.Children = append([]*BTreeNode(nil), t.Root.Children[:mid+1]...)
		right.Children = append([]*BTreeNode(nil), t.Root.Children[mid+1:]...)
		// 分别为左右节点 设置父亲节点
		setParent(left.Children, left)
		setParent(right.Children, right)
	}
	// 产生新的根节点
	newRoot := &BTreeNode{
		Entries:  []*Item{t.Root.Entries[mid]},
		Children: []*BTreeNode{left, right},
	}
	// 更新根节点
	left.Parent = newRoot
	right.Parent = newRoot
	t.Root = newRoot
}

// 分割非根节点
func (t *BTree) splitNonRoot(node *BTreeNode) {
	// 找到mid索引下标
	middle := t.middle()
	// 找到父节点 因为对于当前节点需要进行拆分了
	// 是从最下面开始拆分的
	parent := node.Parent

	// 找到左边部分
	left := &BTreeNode{Entries: append([]*Item(nil), node.Entries[:middle]...), Parent: parent}
	// 找到右边部分
	right := &BTreeNode{Entries: append([]*Item(nil), node.Entries[middle+1:]...), Parent: parent}

	// 如果当前node节点 不是叶子节点的话
	if !t.isLeaf(node) {
		// 对于左边的孩子处理
		left.Children = append([]*BTreeNode(nil), node.Children[:middle+1]...)
		// 对于右边孩子进行处理
		right.Children = append([]*BTreeNode(nil), node.Children[middle+1:]...)
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

	// 继续往上走 开始递归上传即可
	t.split(parent)
}

// 设置父亲节点
func setParent(nodes []*BTreeNode, parent *BTreeNode) {
	for _, node := range nodes {
		node.Parent = parent
	}
}

// 节点个数大于最大个数的话
func (t *BTree) shouldSplit(node *BTreeNode) bool {
	return len(node.Entries) > t.maxEntries()
}
func (t *BTree) Empty() bool {
	return t.size == 0
}
func (t *BTree) middle() int {
	return (t.m - 1) / 2
}
func (t *BTree) maxChildren() int {
	return t.m
}
func (t *BTree) minChildren() int {
	return (t.m + 1) / 2
}
func (t *BTree) maxEntries() int {
	return t.maxChildren() - 1
}
func (t *BTree) minEntries() int {
	return t.minChildren() - 1
}
func (t *BTree) isLeaf(node *BTreeNode) bool {
	return len(node.Children) == 0
}

// 不断地进行二分查找 进行搜索 知道找到最终的叶子节点所在位置
// 节点一定是
func (t *BTree) insertIntoInternal(node *BTreeNode, entry *Item) bool {
	idx, ok := t.search(node, entry.Key)
	if ok {
		node.Entries[idx] = entry
		return false
	}
	return t.insert(node.Children[idx], entry)
}

// 查找操作 简单的二分 靠左查询
func (t *BTree) search(node *BTreeNode, key int) (int, bool) {
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
func (t *BTree) searchRecur(node *BTreeNode, key int) (st *BTreeNode, idx int, ok bool) {
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
func (t *BTree) Get(entry *Item) bool {
	node, idx, ok := t.searchRecur(t.Root, entry.Key)
	entry.Val = node.Entries[idx]
	return ok
}

// Remove 进行移除操作
func (t *BTree) Remove(entry *Item) bool {
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
func (t *BTree) delete(node *BTreeNode, idx int) {
	// 如果是叶子节点 删除叶子节点
	// 否则删除内部节点
	if t.isLeaf(node) {
		t.deleteFromLeaf(node, idx)
	} else {
		t.deleteFromInternal(node, idx)
	}
	// 如果父节点不为空并且自身的节点个数小于最小要求个数
	// 进行向上调整
	if node.Parent != nil && len(node.Entries) < t.minEntries() {
		t.rebalanced(node)
	}
}

// 删除叶子节点
func (t *BTree) deleteFromLeaf(node *BTreeNode, idx int) {
	copy(node.Entries[idx:], node.Entries[idx+1:])
	node.Entries[len(node.Entries)-1] = nil
	node.Entries = node.Entries[:len(node.Entries)-1]
}

// 删除内部节点
func (t *BTree) deleteFromInternal(node *BTreeNode, idx int) {
	// 如果节点的实体
	if len(node.Children[idx].Entries) >= t.minEntries() {
		pred := t.getPredecessor(node, idx)
		node.Entries[idx] = pred
		t.delete(node.Children[idx], len(node.Children[idx].Entries)-1)
	} else if len(node.Children[idx+1].Entries) >= t.minEntries() {
		succ := t.getSuccessor(node, idx)
		node.Entries[idx] = succ
		t.delete(node.Children[idx+1], 0)
	} else {
		t.merge(node, idx)
		t.delete(node.Children[idx], t.minEntries()-1)
	}
}

// 获取前一个节点
func (t *BTree) getPredecessor(node *BTreeNode, idx int) *Item {
	curr := node.Children[idx]
	for !t.isLeaf(curr) {
		curr = curr.Children[len(curr.Children)-1]
	}
	return curr.Entries[len(curr.Entries)-1]
}

// 获取到后一个节点
func (t *BTree) getSuccessor(node *BTreeNode, idx int) *Item {
	curr := node.Children[idx+1]
	for !t.isLeaf(curr) {
		curr = curr.Children[0]
	}
	return curr.Entries[0]
}

// 进行合并操作
func (t *BTree) merge(node *BTreeNode, idx int) {
	//左边+右边
	child := node.Children[idx]
	sibling := node.Children[idx+1]

	//开始合并
	child.Entries = append(child.Entries, node.Entries[idx])
	child.Entries = append(child.Entries, sibling.Entries...)
	// 如果不是叶子节点的话  需要考虑孩子节点
	if !t.isLeaf(child) {
		child.Children = append(child.Children, sibling.Children...)
	}
	//处理一下 实体
	copy(node.Entries[idx:], node.Entries[idx+1:])
	node.Entries = node.Entries[:len(node.Entries)-1]
	//处理孩子节点
	copy(node.Children[idx+1:], node.Children[idx+2:])
	node.Children = node.Children[:len(node.Children)-1]

	if node == t.Root && len(node.Entries) == 0 {
		t.Root = child
		child.Parent = nil
	}
}

// 进行平衡操作
func (t *BTree) rebalanced(node *BTreeNode) {
	//获取父节点
	parent := node.Parent
	if parent == nil {
		return
	}
	// 并且找到该node对应的index
	var idx int
	for i, child := range parent.Children {
		if child == node {
			idx = i
			break
		}
	}
	// 如果左边有节点并且左边的实体可以足够来借
	if idx > 0 && len(parent.Children[idx-1].Entries) > t.minEntries() {
		t.rotateRight(parent, idx-1)
		//	如果右边右节点并且右边的实体是足够来借
	} else if idx < len(parent.Children)-1 && len(parent.Children[idx+1].Entries) > t.minEntries() {
		t.rotateLeft(parent, idx)
		//	在不够借的情况下 仅需要来进行合并操作呢
	} else {
		//左边可以合并的
		if idx > 0 {
			t.merge(parent, idx-1)
			//	右边可以合并
		} else {
			t.merge(parent, idx)
		}
		//合并之后处理一下可能存在的
		//需要继续上溯调整或是已经到达了顶层
		if parent.Parent != nil && len(parent.Entries) < t.minEntries() {
			t.rebalanced(parent)
		} else if parent == t.Root && len(parent.Entries) == 0 {
			t.Root = parent.Children[0]
			t.Root.Parent = nil
		}
	}
}

/*
+--4											+--4,6
     |--  2					删除2之后				 |--  1,3
     |    |--  1								 |--  5
     |    └--  3								 └--  7,8
     └--  6
         |--  5
         └--  7,8

*/
// 左边的节点要进行向右边旋转操作
func (t *BTree) rotateRight(parent *BTreeNode, idx int) {
	child := parent.Children[idx]
	sibling := parent.Children[idx+1]
	//将 父亲节点上的放到右边的第一个位置上
	sibling.Entries = append([]*Item{parent.Entries[idx]}, sibling.Entries...)
	// 将左边的最后一个放在父亲节点上
	parent.Entries[idx] = child.Entries[len(child.Entries)-1]
	// 左边的缩小长度
	child.Entries = child.Entries[:len(child.Entries)-1]

	//如果左边的节点并不是叶子节点 那对于他的孩子节点 还需要处理一下
	if !t.isLeaf(child) {
		//对于右边的孩子进行填充
		sibling.Children = append([]*BTreeNode{child.Children[len(child.Children)-1]}, sibling.Children...)
		//对于最后进行移除
		child.Children = child.Children[:len(child.Children)-1]
	}
}

// 右边的节点进行左边旋转操作 parent, idx
func (t *BTree) rotateLeft(parent *BTreeNode, idx int) {
	child := parent.Children[idx]
	sibling := parent.Children[idx+1]

	child.Entries = append(child.Entries, parent.Entries[idx])
	parent.Entries[idx] = sibling.Entries[0]
	sibling.Entries = sibling.Entries[1:]

	if !t.isLeaf(sibling) {
		child.Children = append(child.Children, sibling.Children[0])
		sibling.Children = sibling.Children[1:]
	}
}

// Inorder 中序遍历
func (t *BTree) Inorder() []*Item {
	if t.Root == nil {
		return []*Item{}
	}
	return inorder(t.Root)
}

func inorder(node *BTreeNode) []*Item {
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
func (t *BTree) LevelOrderTraversal() (ans [][]*Item) {
	return levelOrderTraversal(t.Root)
}
func levelOrderTraversal(root *BTreeNode) (ans [][]*Item) {
	if root == nil {
		return
	}
	queue := []*BTreeNode{root}
	for i := 0; i < len(queue); i++ {
		ans = append(ans, []*Item{})
		var pueue []*BTreeNode
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
