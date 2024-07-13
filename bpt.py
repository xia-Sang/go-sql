class BPlusTreeNode:
    def __init__(self, leaf=False):
        self.leaf = leaf
        self.keys = []
        self.children = []

class BPlusTree:
    def __init__(self, t):
        self.root = BPlusTreeNode(True)
        self.t = t

    def is_leaf(self, node):
        return len(node.children) == 0

    def should_split(self, node):
        return len(node.keys) >= self.t

    def middle(self):
        return (self.t + 1) // 2

    def insert(self, key):
        root = self.root
        if self.should_split(root):
            temp = BPlusTreeNode()
            self.root = temp
            temp.children.append(root)
            self._split_child(temp, 0)
            self._insert_non_full(temp, key)
        else:
            self._insert_non_full(root, key)

    def _insert_non_full(self, node, key):
        if self.is_leaf(node):
            node.keys.append(key)
            node.keys.sort()
        else:
            i = len(node.keys) - 1
            while i >= 0 and key < node.keys[i]:
                i -= 1
            i += 1
            if self.should_split(node.children[i]):
                self._split_child(node, i)
                if key > node.keys[i]:
                    i += 1
            self._insert_non_full(node.children[i], key)

    def _split_child(self, node, i):
        t = self.t
        y = node.children[i]
        z = BPlusTreeNode(y.leaf)
        node.children.insert(i + 1, z)
        node.keys.insert(i, y.keys[self.middle()])
        z.keys = y.keys[t:(2 * t) - 1]
        y.keys = y.keys[0:t - 1]
        if not y.leaf:
            z.children = y.children[t:(2 * t)]
            y.children = y.children[0:t]

    def search(self, key, node=None):
        if node is None:
            node = self.root
        i = 0
        while i < len(node.keys) and key > node.keys[i]:
            i += 1
        if i < len(node.keys) and key == node.keys[i]:
            return True
        if self.is_leaf(node):
            return False
        return self.search(key, node.children[i])

    def level_order_traversal(self):
        if not self.root:
            return

        queue = [self.root]
        while queue:
            current = queue.pop(0)
            print(' '.join(map(str, current.keys)), end=' | ')
            if not self.is_leaf(current):
                queue.extend(current.children)

    def visualize(self):
        levels = []
        self._get_levels(self.root, 0, levels)
        for i, level in enumerate(levels):
            print(f"Level {i}: ", end='')
            for node in level:
                print(f"[{' '.join(map(str, node.keys))}]", end=' ')
            print()

    def _get_levels(self, node, level, levels):
        if len(levels) == level:
            levels.append([])
        levels[level].append(node)
        if not self.is_leaf(node):
            for child in node.children:
                self._get_levels(child, level + 1, levels)

# Usage example
bpt = BPlusTree(3)
for i in range(1,10):
    bpt.insert(i)

bpt.visualize()
bpt.level_order_traversal()
