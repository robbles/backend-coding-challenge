package models

import "strings"

// This is a tree that:
// - stores its children in a map{string -> node}
// - uses single char keys in every node
// - is case-insensitive
type Trie struct {
	edges Edges
	leaf  bool
	value []Location
}

type Edges map[rune]*Trie

func NewTrie() *Trie {
	return &Trie{
		edges: make(Edges),
		leaf:  false,
	}
}

// Insert a key into the tree.
func (tree *Trie) Insert(key string, value Location) {
	var node *Trie = tree

	// iterate through each character in the key, lower-cased
	for _, char := range strings.ToLower(key) {
		// check to see if it exists as a child of the current tree
		leaf, found := node.edges[char]
		if !found {
			// if not found, we need to create the leaf node and attach it
			leaf = NewTrie()
			node.edges[char] = leaf
		}

		// look at the current leaf node in the next iteration
		node = leaf
	}

	// mark the current leaf node as a leaf and store the value
	node.leaf = true
	node.value = append(node.value, value)
}

// Check if a key is present in the tree.
func (tree *Trie) Find(key string) bool {
	var node *Trie = tree
	var found bool

	for _, char := range strings.ToLower(key) {
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
func (tree *Trie) FindMatches(prefix string, limit int) []Location {
	root := tree
	count := 0
	results := []Location{}

	// find the subset of the tree that matches the query,
	// and set that as the current root
	for _, char := range strings.ToLower(prefix) {
		child, found := root.edges[char]
		if !found {
			return results
		}
		root = child
	}

	// breadth first search starting at current root node
	queue := []*Trie{root}
	var node *Trie

	for len(queue) > 0 {
		// dequeue safely (range doesn't allow modifying the original slice)
		node, queue = queue[0], queue[1:]

		// only store leaf nodes as results
		if node.leaf {
			results = append(results, node.value...)
			count += 1

			if limit > 0 && count >= limit {
				break
			}
		}

		// after processing each node, add its children to the end of the queue
		for _, child := range node.edges {
			queue = append(queue, child)
		}
	}

	return results
}
