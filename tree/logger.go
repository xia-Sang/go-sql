package tree

//后续可能会删掉的完善部分
import "fmt"

// Logger 实现打印输出操作
type Logger struct {
	detail bool
}

// levelOrderTraversal 层次遍历函数
func (logger *Logger) levelOrderTraversal(root *BTreeNode) (ans [][]string) {
	if root == nil {
		return
	}
	queue := []*BTreeNode{root}
	for i := 0; i < len(queue); i++ {
		ans = append(ans, []string{})
		var pueue []*BTreeNode
		for j := 0; j < len(queue); j++ {
			node := queue[j]
			curStr := logger.strEntries(node)
			ans[i] = append(ans[i], curStr)
			// 将当前节点的子节点加入队列
			for _, child := range node.Children {
				pueue = append(pueue, child)
			}
		}
		queue = pueue
	}
	return
}

// 递归打印btree
func (logger *Logger) tree(node *BTreeNode, childName string, dsc func(*BTreeNode, bool) string, depth int, prefix string) {
	if depth == 0 {
		fmt.Printf("+--%s\n", dsc(node, logger.detail))
		depth++
	}

	childCount := len(node.Children)
	for idx, child := range node.Children {
		isLastChild := idx == childCount-1

		newPrefix := prefix
		if isLastChild {
			newPrefix += "    "
		} else {
			newPrefix += " |   "
		}

		fmt.Println(prefix, childPrefix(isLastChild), dsc(child, logger.detail))

		logger.tree(child, childName, dsc, depth+1, newPrefix)
	}
}

func childPrefix(isLastChild bool) string {
	if isLastChild {
		return "└-- "
	}
	return "|-- "
}

// strEntries转为字符串类型
func (logger *Logger) strEntries(node *BTreeNode) string {
	s := ""
	for i, keyword := range node.Entries {
		if logger.detail {
			s += fmt.Sprintf("%v", keyword.info())
		} else {
			s += fmt.Sprintf("%v", keyword.Key)
		}
		if i != len(node.Entries)-1 {
			s += ","
		}
	}
	return s
}
func (logger *Logger) PrintTree(node *BTreeNode) {
	fmt.Println("\n************BTree*************")

	dsc := func(node *BTreeNode, detail bool) string {
		return logger.strEntries(node)
	}
	logger.tree(node, "child_nodes", dsc, 0, "    ")

	fmt.Println("******************************")
}

// LevelOrderTraversal 实现层次遍历
func (logger *Logger) LevelOrderTraversal(node *BTreeNode) {
	level := logger.levelOrderTraversal(node)
	for _, data := range level {
		fmt.Println(data)
	}
}