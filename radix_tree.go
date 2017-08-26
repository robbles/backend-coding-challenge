package main

/* This is a Radix Tree that:
   - stores its children in a map{string -> node}
   - uses single char keys in every node
*/
type RadixTree struct {
	edges Edges
	leaf  bool
}

type Edges map[rune]*RadixTree

func NewRadixTree() *RadixTree {
	return &RadixTree{
		edges: make(Edges),
		leaf:  false,
	}
}

// Insert a key into the tree
func (tree *RadixTree) Insert(key string) {
	var node *RadixTree = tree

	// iterate through each character in the key
	for _, char := range key {
		// check to see if it exists as a child of the current tree
		leaf, found := node.edges[char]
		if !found {
			// if not found, we need to create the leaf node and attach it
			leaf = NewRadixTree()
			node.edges[char] = leaf
		}

		// look at the current leaf node in the next iteration
		node = leaf
	}

	// mark the current leaf node as a leaf
	node.leaf = true
}

func (tree *RadixTree) Find(key string) bool {
	var node *RadixTree = tree
	var found bool

	for _, char := range key {
		node, found = node.edges[char]
		if !found {
			return false
		}
		if node.leaf {
			return true
		}
	}

	return false
}

func (tree RadixTree) FindMatches(key string, limit int) []string {
	// TODO: test and implement once RadixTree.Find is complete
	return nil
}
