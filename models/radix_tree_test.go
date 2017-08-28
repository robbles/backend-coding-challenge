package models

import (
	"reflect"
	"testing"
)

func makeTree(char rune, child *RadixTree) *RadixTree {
	tree := NewRadixTree()
	tree.edges[char] = child
	return tree
}

func makeLeaf(value string) *RadixTree {
	tree := NewRadixTree()
	tree.leaf = true
	tree.value = value
	return tree
}

func TestRadixTree_Find(t *testing.T) {
	tests := map[string]struct {
		tree     *RadixTree
		key      string
		expected bool
	}{
		"empty tree": {
			NewRadixTree(),
			"nope",
			false,
		},
		"missing key": {
			makeTree('a', makeLeaf("a")),
			"nope",
			false,
		},
		"single character": {
			makeTree('a', makeLeaf("a")),
			"a",
			true,
		},
		"multiple characters": {
			makeTree('a', makeTree('b', makeTree('c', makeLeaf("abc")))),
			"abc",
			true,
		},
		"multiple character subset": {
			makeTree('a', makeTree('b', makeTree('c', makeLeaf("abc")))),
			"ab",
			false,
		},
		"multiple character superset": {
			makeTree('a', makeTree('b', makeTree('c', makeLeaf("abc")))),
			"abcd",
			false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if actual := tt.tree.Find(tt.key); actual != tt.expected {
				t.Errorf("%v != %v", actual, tt.expected)
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
		"single character empty tree": {
			NewRadixTree(),
			"a",
			makeTree('a', makeLeaf("a")),
		},
		"empty key empty tree": {
			NewRadixTree(),
			"",
			makeLeaf(""),
		},
		"multiple characters empty tree": {
			NewRadixTree(),
			"abc",
			makeTree('a', makeTree('b', makeTree('c', makeLeaf("abc")))),
		},
		"multiple characters non-empty tree": {
			makeTree('a', makeTree('b', makeTree('c', makeLeaf("abc")))),
			"abd",
			makeTree('a', makeTree('b', &RadixTree{
				edges: Edges{
					'c': makeLeaf("abc"),
					'd': makeLeaf("abd"),
				},
			})),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tree := tt.before
			if tree.Insert(tt.key); !reflect.DeepEqual(tree, tt.after) {
				t.Errorf("\n%#v\n!=\n%#v", tree, tt.after)
			}
		})
	}
}

func TestRadixTree_FindMatches(t *testing.T) {
	tests := map[string]struct {
		tree     *RadixTree
		key      string
		limit    int
		expected []string
	}{
		"empty tree": {
			NewRadixTree(),
			"nope",
			10,
			[]string{},
		},
		"missing key": {
			makeTree('a', makeLeaf("a")),
			"nope",
			10,
			[]string{},
		},
		"exact match": {
			makeTree('a', makeLeaf("a")),
			"a",
			10,
			[]string{"a"},
		},
		"multiple matches": {
			makeTree('a', makeTree('b', &RadixTree{
				edges: Edges{
					'c': makeLeaf("abc"),
					'd': makeLeaf("abd"),
				},
			})),
			"ab",
			10,
			[]string{"abc", "abd"},
		},
		"multiple matches limit shortest first": {
			makeTree('a', makeTree('b', &RadixTree{
				edges: Edges{
					'c': makeLeaf("abc"),
					'd': makeTree('e', makeLeaf("abde")),
				},
			})),
			"ab",
			1,
			[]string{"abc"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if actual := tt.tree.FindMatches(tt.key, tt.limit); !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("%#v != %#v", actual, tt.expected)
			}
		})
	}
}
