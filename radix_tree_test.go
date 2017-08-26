package main

import (
	"reflect"
	"testing"
)

func makeTree(char rune, child *RadixTree) *RadixTree {
	tree := NewRadixTree()
	tree.edges[char] = child
	return tree
}

func makeLeaf() *RadixTree {
	tree := NewRadixTree()
	tree.leaf = true
	return tree
}

func TestRadixTree_Find(t *testing.T) {
	tests := map[string]struct {
		tree *RadixTree
		key  string
		want bool
	}{
		"empty tree": {
			NewRadixTree(),
			"nope",
			false,
		},
		"missing key": {
			makeTree('a', makeLeaf()),
			"nope",
			false,
		},
		"single character": {
			makeTree('a', makeLeaf()),
			"a",
			true,
		},
		"multiple characters": {
			makeTree('a', makeTree('b', makeTree('c', makeLeaf()))),
			"abc",
			true,
		},
		"multiple character subset": {
			makeTree('a', makeTree('b', makeTree('c', makeLeaf()))),
			"ab",
			false,
		},
		"multiple character superset": {
			makeTree('a', makeTree('b', makeTree('c', makeLeaf()))),
			"abcd",
			false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.tree.Find(tt.key); got != tt.want {
				t.Errorf("RadixTree.Find() %v != %v", got, tt.want)
			}
		})
	}
}

func TestRadixTree_Insert(t *testing.T) {
	tests := map[string]struct {
		before *RadixTree
		key    string
		after  *RadixTree
	}{
		"empty key empty tree": {
			NewRadixTree(),
			"a",
			makeTree('a', makeLeaf()),
		},
		"single character empty tree": {
			NewRadixTree(),
			"",
			makeLeaf(),
		},
		"multiple characters empty tree": {
			NewRadixTree(),
			"abc",
			makeTree('a', makeTree('b', makeTree('c', makeLeaf()))),
		},
		"multiple characters non-empty tree": {
			makeTree('a', makeTree('b', makeTree('c', makeLeaf()))),
			"abd",
			makeTree('a', makeTree('b', &RadixTree{
				edges: Edges{
					'c': makeLeaf(),
					'd': makeLeaf(),
				},
			})),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.before.Insert(tt.key); !reflect.DeepEqual(tt.before, tt.after) {
				t.Errorf("RadixTree.Insert()\n%v\n!=\n%v", tt.before, tt.after)
			}
		})
	}
}
