package bptree

import (
	"fmt"
	"testing"
)

func TestNewTree(t *testing.T) {
	tree := NewBTree(3)
	for i := range 10 {
		tree.Put(i, fmt.Sprintf("数据%d", i))
	}
	printTree(tree.Root)
	tree.Put(64, fmt.Sprintf("数据%d", 64))
	for _, i := range []int{20, 30, 50, 52, 60, 62, 64} {
		val, ok := tree.Get(i)
		t.Log(i, val, ok)
	}
	printTree(tree.Root)
	tree.Remove(52)
	printTree(tree.Root)
}
