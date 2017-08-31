package models

/* This is a Radix Tree that:
   - stores its children in a map{string -> node}
   - uses single char keys in every node
*/
type RadixTree struct {
	edges Edges

	// TODO: should this field be removed since only leaf nodes have values?
	leaf bool

	// TODO: does this need to become an interface{}?
	// or is it better to have RadixTree know about city objects?
	value string
}

type Edges map[rune]*RadixTree

func NewRadixTree() *RadixTree {
	return &RadixTree{
		edges: make(Edges),
		leaf:  false,
	}
}

// Insert a key into the tree.
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

	// mark the current leaf node as a leaf and store the value
	node.leaf = true
	node.value = key
}

// Check if a key is present in the tree.
func (tree *RadixTree) Find(key string) bool {
	var node *RadixTree = tree
	var found bool

	for _, char := range key {
		node, found = node.edges[char]
		if !found {
			return false
		}
	}

	if node != nil && node.leaf {
		return true
	}

	return false
}

// Find <limit> matches with the given <prefix>.
func (tree *RadixTree) FindMatches(prefix string, limit int) []string {
	root := tree
	count := 0
	results := []string{}

	// find the subset of the tree that matches the query,
	// and set that as the current root
	for _, char := range prefix {
		child, found := root.edges[char]
		if !found {
			return results
		}
		root = child
	}

	// breadth first search starting at current root node
	queue := []*RadixTree{root}
	var node *RadixTree

	for len(queue) > 0 {
		// dequeue safely (range doesn't allow modifying the original slice)
		node, queue = queue[0], queue[1:]

		// only store leaf nodes as results
		if node.leaf {
			results = append(results, node.value)
			count += 1

			if limit > 0 && count >= limit {
				break
			}
		}

		// after processing each node, add its children to the end of the queue
		for _, value := range node.edges {
			queue = append(queue, value)
		}
	}

	// TODO: sort results and assign scores

	return results
}
