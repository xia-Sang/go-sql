package tree

import (
	"fmt"
	"github.com/xia-Sang/go-sql/util"
	"testing"
)

func TestNewTree(t *testing.T) {
	tree := NewBTree(3)
	for i := range 9 {
		tree.Put(&Item{i - 9, fmt.Sprintf("数据%d", i)})
	}
	log := Logger{}
	log.PrintTree(tree.Root)
	log.LevelOrderTraversal(tree.Root)
	tree.Remove(&Item{Key: 2})
	log.PrintTree(tree.Root)
	log.levelOrderTraversal(tree.Root)
	tree.Put(&Item{Key: 1, Val: "dkebn"})
	log.PrintTree(tree.Root)
	log.LevelOrderTraversal(tree.Root)
	order := tree.Inorder()
	t.Log(order)
	for _, o := range order {
		t.Log(o.info())
	}
	t.Log(levelOrderTraversal(tree.Root))
}
func TestNewTree1(t *testing.T) {
	tree := NewBTree(3)
	for _, i := range util.GenerateRandomNumbers(12, 0, 100) {
		tree.Put(&Item{i, fmt.Sprintf("数据%d", i)})
	}
	log := Logger{}
	log.PrintTree(tree.Root)
	log.LevelOrderTraversal(tree.Root)
	tree.Remove(&Item{Key: 2})
	log.PrintTree(tree.Root)
	log.levelOrderTraversal(tree.Root)
	tree.Put(&Item{Key: 1, Val: "dkebn"})
	log.PrintTree(tree.Root)
	log.LevelOrderTraversal(tree.Root)
	order := tree.Inorder()
	t.Log(order)
	for _, o := range order {
		t.Log(o.info())
	}

	t.Log(levelOrderTraversal(tree.Root))
	t.Log(util.GenerateRandomNumbers(12, 0, 100))
}