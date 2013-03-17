package movegen

import (

)

type Tree struct {
	State GameState
	Children []*Tree
}

func AddChild(parent *Tree, child *Tree) {
	parent.Children = append(parent.Children, child)
}

func AddChildren(parent *Tree, children...*Tree) {
	parent.Children = append(parent.Children, children...)
}

func GameStatesToTrees(states []GameState) []*Tree {
	trees := make([]*Tree, 0)
	for _, val := range states {
		tree := Tree { val, make([]*Tree, 0) }
		trees = append(trees, &tree)
	}
	return trees
}
