package tree

import (
	"fmt"
	"testing"

	"github.com/xia-Sang/go-sql/util"
)

func TestNewTree(t *testing.T) {
	tree := NewBPTree(3)

	log := Logger{true}
	tree.Put(&Item{1, fmt.Sprintf("数据%d", 1), false})
	log.PrintTree(tree.Root)
	log.LevelOrderTraversal(tree.Root)

	tree.Put(&Item{Key: 2})
	log.PrintTree(tree.Root)
	log.LevelOrderTraversal(tree.Root)

	tree.Put(&Item{Key: 3})
	log.PrintTree(tree.Root)
	log.LevelOrderTraversal(tree.Root)
	log.PrintTree(tree.Root)
	log.LevelOrderTraversal(tree.Root)
	tree.Put(&Item{Key: 4})
	log.PrintTree(tree.Root)
	log.LevelOrderTraversal(tree.Root)

	tree.Put(&Item{Key: 5})
	log.PrintTree(tree.Root)
	log.LevelOrderTraversal(tree.Root)

	// tree.Put(&Item{Key: 6})
	// log.PrintTree(tree.Root)
	// log.LevelOrderTraversal(tree.Root)
}
func TestNewTree1(t *testing.T) {
	tree := NewBPTree(3)
	for _, i := range util.GenerateRandomNumbers(12, 0, 100) {
		tree.Put(&Item{i, fmt.Sprintf("数据%d", i), false})
	}
	log := Logger{}
	log.PrintTree(tree.Root)
	log.LevelOrderTraversal(tree.Root)
	// tree.Remove(&Item{Key: 2})
	// log.PrintTree(tree.Root)
	// log.levelOrderTraversal(tree.Root)
	// tree.Put(&Item{Key: 1, Val: "dkebn"})
	// log.PrintTree(tree.Root)
	// log.LevelOrderTraversal(tree.Root)
	// order := tree.Inorder()
	// t.Log(order)
	// for _, o := range order {
	// 	t.Log(o.info())
	// }

	// t.Log(levelOrderTraversal(tree.Root))
	// t.Log(util.GenerateRandomNumbers(12, 0, 100))
}
