package tree

// import (
// 	"fmt"
// )

// // BPTree 存储节点
// type BPTree struct {
// 	Root      *BPTreeNode // 根节点
// 	size      int         // 统计规模
// 	m         int         // 度
// 	splitShow bool        // 节点分列是否显示 默认不显示
// }

// // BPTreeNode 存储node 节点
// type BPTreeNode struct {
// 	Parent   *BPTreeNode   // 父节点
// 	Entries  []*Item       // 存储实体
// 	Children []*BPTreeNode // 子节点 列表
// }

// // Item 存储的实体数据
// type Item struct {
// 	Key int
// 	Val interface{}
// }

// func (i *Item) info() string {
// 	return fmt.Sprintf("(%d:%v)", i.Key, i.Val)
// }

// // NewBPTree 构建新树
// func NewBPTree(order int) *BPTree {
// 	if order < 3 {
// 		order = 3
// 	}
// 	return &BPTree{m: order}
// }

// // Put 进行put操作 存储数据
// func (t *BPTree) Put(entry *Item) {
// 	// 如果根节点为空，直接结束
// 	if t.Root == nil {
// 		t.Root = &BPTreeNode{Entries: []*Item{entry}}
// 		t.size++
// 		return
// 	}
// 	// 根节点不为空，进行插入操作，插入成功更新 size
// 	if t.insert(t.Root, entry) {
// 		t.size++
// 	}
// }

// // 进行插入操作
// func (t *BPTree) insert(node *BPTreeNode, entry *Item) bool {
// 	// 如果是叶子节点，进行叶子节点插入
// 	if t.isLeaf(node) {
// 		return t.insertIntoLeaf(node, entry)
// 	}
// 	// 不是叶子节点，继续向下查找
// 	return t.insertIntoInternal(node, entry)
// }

// // 不断地进行二分查找 进行搜索 知道找到最终的叶子节点所在位置
// // 节点一定是

// func (t *BPTree) insertIntoInternal(node *BPTreeNode, entry *Item) bool {
// 	idx, ok := t.search(node, entry.Key)
// 	if ok {
// 		// node.Entries[idx] = entry
// 		// return false
// 		setNewEntry(node.Entries, idx, entry)
// 		return false
// 	}
// 	return t.insert(node.Children[idx], entry)
// }

// // 设置新条目
// func setNewEntry(entries []*Item, index int, entry *Item) {
// 	entries[index].Key = entry.Key
// 	entries[index].Val = entry.Val
// }

// // 叶子节点的插入
// // 找到对应的叶子节点
// // 返回 false 表示修改，true 表示插入
// func (t *BPTree) insertIntoLeaf(node *BPTreeNode, entry *Item) bool {
// 	// 进行节点查找
// 	idx, ok := t.search(node, entry.Key)
// 	if ok {
// 		setNewEntry(node.Entries, idx, entry)
// 		return false
// 	}
// 	// 开辟新的空间，多一个 nil
// 	node.Entries = append(node.Entries, nil)
// 	// 复制数据
// 	copy(node.Entries[idx+1:], node.Entries[idx:])
// 	// 在对应位置插入数据
// 	node.Entries[idx] = entry
// 	// 对当前节点进行 split 操作
// 	t.split(node)
// 	return true
// }

// // 进行 split 操作
// func (t *BPTree) split(node *BPTreeNode) {
// 	// 检查是否需要 split 操作
// 	if !t.shouldSplit(node) {
// 		return
// 	}
// 	if t.splitShow {
// 		fmt.Println("node:", node, "start split")
// 	}
// 	// 根节点特殊处理
// 	if node == t.Root {
// 		t.splitRoot()
// 		return
// 	}
// 	if t.isLeaf(node) {
// 		t.splitLeaf(node)
// 	} else {
// 		// 非叶子节点处理
// 		t.splitNonLeaf(node)
// 	}
// }

// // 分割根节点
// func (t *BPTree) splitRoot() {
// 	// 找到中间索引
// 	mid := t.middle()
// 	// 左边节点
// 	left := &BPTreeNode{
// 		Entries: append([]*Item(nil), t.Root.Entries[:mid]...),
// 	}
// 	// 右边节点
// 	right := &BPTreeNode{
// 		Entries: append([]*Item(nil), t.Root.Entries[mid:]...),
// 	}
// 	// 如果根节点不是叶子节点
// 	if !t.isLeaf(t.Root) {
// 		left.Children = append([]*BPTreeNode(nil), t.Root.Children[:mid+1]...)
// 		right.Children = append([]*BPTreeNode(nil), t.Root.Children[mid+1:]...)
// 		// 设置父节点
// 		setParent(left.Children, left)
// 		setParent(right.Children, right)
// 	}
// 	// 产生新的根节点
// 	newRoot := &BPTreeNode{
// 		Entries:  []*Item{t.Root.Entries[mid]},
// 		Children: []*BPTreeNode{left, right},
// 	}
// 	// 更新根节点
// 	left.Parent = newRoot
// 	right.Parent = newRoot
// 	t.Root = newRoot
// }

// // 设置父节点
// func setParent(nodes []*BPTreeNode, parent *BPTreeNode) {
// 	for _, node := range nodes {
// 		node.Parent = parent
// 	}
// }

// // 查找节点
// func (t *BPTree) search(node *BPTreeNode, key int) (int, bool) {
// 	for i, entry := range node.Entries {
// 		if entry.Key == key {
// 			return i, true
// 		} else if entry.Key > key {
// 			return i, false
// 		}
// 	}
// 	return len(node.Entries), false
// }

// // 分割叶子节点
// func (t *BPTree) splitLeaf(node *BPTreeNode) {
// 	// 找到中间索引
// 	mid := t.middle()
// 	// 新建右节点
// 	right := &BPTreeNode{
// 		Entries: append([]*Item(nil), node.Entries[mid+1:]...),
// 		Parent:  node.Parent,
// 	}
// 	// 更新当前节点
// 	node.Entries = node.Entries[:mid]
// 	// 插入到父节点
// 	t.insertIntoParent(node, right, node.Entries[mid])
// }

// // 分割非叶子节点
// func (t *BPTree) splitNonLeaf(node *BPTreeNode) {
// 	// 找到中间索引
// 	mid := t.middle()
// 	// 新建右节点
// 	right := &BPTreeNode{
// 		Entries:  append([]*Item(nil), node.Entries[mid:]...),
// 		Children: append([]*BPTreeNode(nil), node.Children[mid+1:]...),
// 		Parent:   node.Parent,
// 	}
// 	// 更新当前节点
// 	node.Entries = node.Entries[:mid]
// 	node.Children = node.Children[:mid+1]
// 	// 设置父节点
// 	setParent(right.Children, right)
// 	// 插入到父节点
// 	t.insertIntoParent(node, right, node.Entries[mid])
// }

// // 插入到父节点
// func (t *BPTree) insertIntoParent(left, right *BPTreeNode, entry *Item) {
// 	if left.Parent == nil {
// 		// 创建新的根节点
// 		t.Root = &BPTreeNode{
// 			Entries:  []*Item{entry},
// 			Children: []*BPTreeNode{left, right},
// 		}
// 		left.Parent = t.Root
// 		right.Parent = t.Root
// 		return
// 	}
// 	parent := left.Parent
// 	// 查找插入位置
// 	idx, _ := t.search(parent, entry.Key)
// 	// 插入父节点
// 	parent.Entries = append(parent.Entries, nil)
// 	copy(parent.Entries[idx+1:], parent.Entries[idx:])
// 	parent.Entries[idx] = entry
// 	parent.Children = append(parent.Children, nil)
// 	copy(parent.Children[idx+2:], parent.Children[idx+1:])
// 	parent.Children[idx+1] = right
// 	right.Parent = parent
// 	// 继续分裂父节点
// 	t.split(parent)
// }

// // 是否为叶子节点
// func (t *BPTree) isLeaf(node *BPTreeNode) bool {
// 	return len(node.Children) == 0
// }

// // 是否需要分裂
// func (t *BPTree) shouldSplit(node *BPTreeNode) bool {
// 	return len(node.Entries) >= t.m
// }

// // 获取中间索引
// func (t *BPTree) middle() int {
// 	return (t.m - 1) / 2
// }
