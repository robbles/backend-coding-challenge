package models

import (
	"reflect"
	"sort"
	"testing"
)

func makeTree(args ...interface{}) *Trie {
	tree := NewTrie()
	for len(args) >= 2 {
		char, child, rest := args[0], args[1], args[2:]
		tree.edges[char.(rune)] = child.(*Trie)
		args = rest
	}
	return tree
}

func makeLeaf(value string) *Trie {
	tree := NewTrie()
	tree.leaf = true
	tree.value = value
	return tree
}

func TestTrie_Find(t *testing.T) {
	tests := map[string]struct {
		tree     *Trie
		key      string
		expected bool
	}{
		"empty tree": {
			NewTrie(),
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

func TestTrie_Insert(t *testing.T) {
	tests := map[string]struct {
		before *Trie
		key    string
		after  *Trie
	}{
		"single character empty tree": {
			NewTrie(),
			"a",
			makeTree('a', makeLeaf("a")),
		},
		"empty key empty tree": {
			NewTrie(),
			"",
			makeLeaf(""),
		},
		"multiple characters empty tree": {
			NewTrie(),
			"abc",
			makeTree('a', makeTree('b', makeTree('c', makeLeaf("abc")))),
		},
		"multiple characters non-empty tree": {
			makeTree('a', makeTree('b', makeTree('c', makeLeaf("abc")))),
			"abd",
			makeTree(
				'a', makeTree(
					'b', makeTree(
						'c', makeLeaf("abc"),
						'd', makeLeaf("abd"),
					),
				),
			),
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

func TestTrie_FindMatches(t *testing.T) {
	tests := map[string]struct {
		tree     *Trie
		key      string
		limit    int
		expected []string
	}{
		"empty tree": {
			NewTrie(),
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
			makeTree(
				'a', makeTree(
					'b', makeTree(
						'c', makeLeaf("abc"),
						'd', makeLeaf("abd"),
					),
				),
			),
			"ab",
			10,
			[]string{"abc", "abd"},
		},
		"multiple matches limit returns shortest first": {
			makeTree('a', makeTree('b', makeTree(
				'c', makeLeaf("abc"),
				'd', makeTree('e', makeLeaf("abde")),
			))),
			"ab",
			1,
			[]string{"abc"},
		},
		"limit < 0 means no limit": {
			makeTree(
				'a', makeTree(
					'b', makeTree(
						'c', makeLeaf("abc"),
						'd', makeLeaf("abd"),
					),
				),
			),
			"ab",
			-1,
			[]string{"abc", "abd"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.tree.FindMatches(tt.key, tt.limit)
			sort.Strings(actual)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("%#v != %#v", actual, tt.expected)
			}
		})
	}
}
