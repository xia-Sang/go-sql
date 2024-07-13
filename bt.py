# 实现数据存储的item
class Item:
    def __init__(self, key, val,deleted=False):
        self.key = key
        self.val = val
        self.deleted=deleted
    def info(self):
        return f"({self.key}:{self.val}:{self.deleted})"

# 实现存储的Node
class BTreeNode:
    def __init__(self, parent=None):
        self.parent = parent
        self.entries = []
        self.children = []
        
# 实现BTree
class BTree:
    # 初始化
    def __init__(self, order):
        self.root = None
        self.size = 0
        self.m = max(order, 3)
        self.split_show = False

    # 进行put操作
    def put(self, entry):
        if self.root is None:
            self.root = BTreeNode()
            self.root.entries.append(entry)
            self.size += 1
            return

        if self._insert(self.root, entry):
            self.size += 1

    # insert插入 注意最终都是要插入到叶子节点的
    def _insert(self, node, entry):
        if self._is_leaf(node):
            return self._insert_into_leaf(node, entry)
        return self._insert_into_internal(node, entry)

    # 插入到叶子节点
    def _insert_into_leaf(self, node, entry):
        idx, ok = self._search(node, entry.key)
        if ok:
            node.entries[idx] = entry
            return False

        node.entries.insert(idx, entry)
        self._split(node)
        return True
    
    #  插入到非叶子部分 进行查找
    def _insert_into_internal(self, node, entry):
        idx, ok = self._search(node, entry.key)
        if ok:
            node.entries[idx] = entry
            return False
        return self._insert(node.children[idx], entry)
    
    # 插入之后可能会产生分裂的可能 一般会有两种
    # 根节点分裂和非根节点分裂
    def _split(self, node):
        if not self._should_split(node):
            return
        if self.split_show:
            print("node:", [e.info() for e in node.entries], "start split")

        if node == self.root:
            self._split_root()
        else:
            self._split_non_root(node)
    
    # 根节点分裂
    def _split_root(self):
        mid = self._middle()
        left = BTreeNode()
        left.entries = self.root.entries[:mid]
        right = BTreeNode()
        right.entries = self.root.entries[mid+1:]
        if not self._is_leaf(self.root):
            left.children = self.root.children[:mid+1]
            right.children = self.root.children[mid+1:]
            self._set_parent(left.children, left)
            self._set_parent(right.children, right)
        new_root = BTreeNode()
        new_root.entries.append(self.root.entries[mid])
        new_root.children = [left, right]
        left.parent = new_root
        right.parent = new_root
        self.root = new_root

    # 非根节点分裂
    def _split_non_root(self, node):
        middle = self._middle()
        parent = node.parent
        left = BTreeNode(parent=parent)
        right = BTreeNode(parent=parent)
        left.entries = node.entries[:middle]
        right.entries = node.entries[middle+1:]

        if not self._is_leaf(node):
            left.children = node.children[:middle+1]
            right.children = node.children[middle+1:]
            self._set_parent(left.children, left)
            self._set_parent(right.children, right)

        insert_position, _ = self._search(parent, node.entries[middle].key)
        parent.entries.insert(insert_position, node.entries[middle])
        parent.children[insert_position] = left
        parent.children.insert(insert_position + 1, right)

        self._split(parent)

    def _set_parent(self, nodes, parent):
        for node in nodes:
            node.parent = parent

    def _should_split(self, node):
        return len(node.entries) > self._max_entries()

    def empty(self):
        return self.size == 0

    def _middle(self):
        return (self.m - 1) // 2

    def _max_children(self):
        return self.m

    def _min_children(self):
        return (self.m + 1) // 2

    def _max_entries(self):
        return self._max_children() - 1

    def _min_entries(self):
        return self._min_children() - 1

    def _is_leaf(self, node:BTreeNode):
        return len(node.children) == 0

    # 实现搜索功能
    def _search(self, node, key):
        left = 0
        right = len(node.entries) - 1
        ans = -1

        while left <= right:
            mid = (left + right) // 2
            if node.entries[mid].key == key:
                return mid, True
            elif node.entries[mid].key < key:
                ans = mid
                left = mid + 1
            else:
                right = mid - 1
        return ans + 1, False

    # 递归搜索
    def _search_recur(self, node, key):
        if self.empty():
            return None, -1, False
        while True:
            idx, ok = self._search(node, key)
            if ok:
                return node, idx, True
            if self._is_leaf(node):
                return None, -1, False
            node = node.children[idx]

    # get查询数据
    def get(self, entry):
        node, idx, ok = self._search_recur(self.root, entry.key)
        if ok:
            entry.val = node.entries[idx].val
        return ok

    # 删除数据 这个稍微比较简单 相当于b+树来说的话
    # 这里的实现不需要太过于复杂
    def remove(self, entry):
        node, idx, found = self._search_recur(self.root, entry)
        if found:
            self._delete(node, idx)
        return found

    def _delete(self, node, idx):
        if self._is_leaf(node):
            self._delete_from_leaf(node, idx)
        else:
            self._delete_from_internal(node, idx)
    

    def remove(self, entry):
        node, idx, ok = self._search_recur(self.root, entry.key)
        if ok:
            self._delete(node, idx)
            self.size -= 1
        return ok

    # 进行懒删除
    def _delete(self, node:BTreeNode, idx):
        node.entries[idx].deleted=True
    

# 实现遍历和打印功能
class BTreeLogger:
    def __init__(self):
        self.detail = False

    def _tree(self, node, child_name, dsc, depth=0, prefix="   "):
        if depth == 0:
            print(f"+--{dsc(node, self.detail)}")
            depth += 1

        child_count = len(node.children)
        for idx, child in enumerate(node.children):
            is_last_child = idx == child_count - 1
            new_prefix = prefix
            if is_last_child:
                new_prefix += "    "
            else:
                new_prefix += "|  "

            print(f"{prefix}{self._child_prefix(is_last_child)}{dsc(child, self.detail)}")
            self._tree(child, child_name, dsc, depth + 1, new_prefix)
    
    def _child_prefix(self, is_last_child):
        if is_last_child:
            return "└-- "
        return "|--"
    
    def inorder(self,node):
        if node is None:
            print("[]")
        ls=self._inorder(node)
        st="["
        for k,v in enumerate(ls):
            st+=v.info()
            if k!=len(ls)-1:
                st+=","
        print(st+"]")

    def _inorder(self, node):
        result = []
        for i in range(len(node.entries)):
            if i < len(node.children):
                result.extend(self._inorder(node.children[i]))
            result.append(node.entries[i])
        if len(node.children) > len(node.entries):
            result.extend(self._inorder(node.children[-1]))
        return result

    def _level_order_traversal(self, root):
        if not root:
            return []
        queue = [root]
        ans = []
        while queue:
            level = []
            next_queue = []
            for node in queue:
                level.extend(node.entries)
                next_queue.extend(node.children)
            ans.append(level)
            queue = next_queue
        return ans
    def level_order_traversal(self, root):
        for data in self._level_order_traversal(root):
            for v in data:
                print(v.info(),sep=",",end=" ")
            print()

    def tree(self, root):
        self._tree(root, "child", lambda node, detail: ",".join([item.info() for item in node.entries]))

def test_btree():
    # 创建一个阶数为3的B树
    btree = BTree(order=3)
    logger = BTreeLogger()

    for i in range(1,10):
        btree.put(Item(i,None))
    # 打印树结构
    print("BTree structure after insertion:")
    logger.tree(btree.root)
    print("*"*20)
    
    btree.remove(Item(2,None))
    logger.tree(btree.root)
    print("*"*20)
    
    btree.put(Item(1,"1"))
    logger.tree(btree.root)
    print("All tests passed!")
     # 创建日志对象
    
    btree.remove(Item(6,None))
    # 打印 B 树结构
    logger.tree(btree.root)
    logger.level_order_traversal(btree.root)
    logger.inorder(btree.root)
# 运行测试
test_btree()
