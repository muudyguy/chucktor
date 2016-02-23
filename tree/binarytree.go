package tree
import "fmt"

// A Tree is a binary tree with integer values.
type InnerBinaryTree struct {
	Left  *InnerBinaryTree
	Value Comparable
	Right *InnerBinaryTree
	Root *InnerBinaryTree
	Parent *InnerBinaryTree
}

func adder(treeNode *InnerBinaryTree, value Comparable) {

	if (value).Bigger(treeNode.Value) {
		if treeNode.Right == nil {
			newNode := &InnerBinaryTree{Left:nil, Value:value, Right:nil, Root:treeNode.Root, Parent:treeNode}
			treeNode.Right = newNode
		}
		adder(treeNode.Right, value)
	}

	if (treeNode.Value).Bigger(value) {
		if treeNode.Left == nil {
			newNode := &InnerBinaryTree{Left:nil, Value:value, Right:nil, Root:treeNode.Root, Parent:treeNode}
			treeNode.Left = newNode
		}
		adder(treeNode.Left, value)
	}
}

func (treeNode *InnerBinaryTree) Add(value Comparable) error {

	if treeNode != treeNode.Root {
		return fmt.Errorf("Cannot add to a non-root node of a tree")
	}

	if treeNode.Value == nil {
		treeNode.Value = value
		return nil
	}

	adder(treeNode, value)
	return nil
}


func searcher(treeNode *InnerBinaryTree, value Comparable, result **InnerBinaryTree) {
	if *result != nil {
		return
	}

	if (value).Equals(treeNode.Value) {
		*result = treeNode
	}

	if (value).Bigger(treeNode.Value) {
		searcher(treeNode.Right, value, result)
	}

	if (treeNode.Value).Bigger(value) {
		searcher(treeNode.Left, value, result)
	}
}


func (treeNode InnerBinaryTree) Search(value Comparable) *InnerBinaryTree {
	var resultHolder **InnerBinaryTree = new(*InnerBinaryTree)
	searcher(&treeNode, value, resultHolder)
	if resultHolder == nil {
		return nil
	} else {
		return *resultHolder
	}
}


type BinaryTree struct {
	root *InnerBinaryTree
}

func (binaryTree BinaryTree) Search(value interface{}) interface{} {
	comparable := (value).(Comparable)
	return binaryTree.root.Search(comparable).Value
}

func (binaryTree BinaryTree) Add(value interface{}) {
	comparable := (value).(Comparable)
	binaryTree.root.Add(comparable)
}

func NewBinaryTree() BinaryTree {
	bt := BinaryTree{}
	bt.root = new(InnerBinaryTree)
	bt.root.Root = bt.root
	return bt
}