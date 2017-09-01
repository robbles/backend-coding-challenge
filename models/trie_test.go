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

func makeLeaf(value City) *Trie {
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
			makeTree('a', makeLeaf(City{Name: "a"})),
			"nope",
			false,
		},
		"single character": {
			makeTree('a', makeLeaf(City{Name: "a"})),
			"a",
			true,
		},
		"multiple characters": {
			makeTree('a', makeTree('b', makeTree('c', makeLeaf(City{Name: "abc"})))),
			"abc",
			true,
		},
		"multiple character subset": {
			makeTree('a', makeTree('b', makeTree('c', makeLeaf(City{Name: "abc"})))),
			"ab",
			false,
		},
		"multiple character superset": {
			makeTree('a', makeTree('b', makeTree('c', makeLeaf(City{Name: "abc"})))),
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
		value  City
		after  *Trie
	}{
		"single character empty tree": {
			NewTrie(),
			"a",
			City{Name: "a"},
			makeTree('a', makeLeaf(City{Name: "a"})),
		},
		"empty key empty tree": {
			NewTrie(),
			"",
			City{Name: ""},
			makeLeaf(City{Name: ""}),
		},
		"multiple characters empty tree": {
			NewTrie(),
			"abc",
			City{Name: "abc"},
			makeTree('a', makeTree('b', makeTree('c', makeLeaf(City{Name: "abc"})))),
		},
		"multiple characters non-empty tree": {
			makeTree('a', makeTree('b', makeTree('c', makeLeaf(City{Name: "abc"})))),
			"abd",
			City{Name: "abd"},
			makeTree(
				'a', makeTree(
					'b', makeTree(
						'c', makeLeaf(City{Name: "abc"}),
						'd', makeLeaf(City{Name: "abd"}),
					),
				),
			),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tree := tt.before
			if tree.Insert(tt.key, tt.value); !reflect.DeepEqual(tree, tt.after) {
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
		expected []City
	}{
		"empty tree": {
			NewTrie(),
			"nope",
			10,
			[]City{},
		},
		"missing key": {
			makeTree('a', makeLeaf(City{Name: "a"})),
			"nope",
			10,
			[]City{},
		},
		"exact match": {
			makeTree('a', makeLeaf(City{Name: "a"})),
			"a",
			10,
			[]City{{Name: "a"}},
		},
		"multiple matches": {
			makeTree(
				'a', makeTree(
					'b', makeTree(
						'c', makeLeaf(City{Name: "abc"}),
						'd', makeLeaf(City{Name: "abd"}),
					),
				),
			),
			"ab",
			10,
			[]City{{Name: "abc"}, {Name: "abd"}},
		},
		"multiple matches limit returns shortest first": {
			makeTree('a', makeTree('b', makeTree(
				'c', makeLeaf(City{Name: "abc"}),
				'd', makeTree('e', makeLeaf(City{Name: "abde"})),
			))),
			"ab",
			1,
			[]City{{Name: "abc"}},
		},
		"limit < 0 means no limit": {
			makeTree(
				'a', makeTree(
					'b', makeTree(
						'c', makeLeaf(City{Name: "abc"}),
						'd', makeLeaf(City{Name: "abd"}),
					),
				),
			),
			"ab",
			-1,
			[]City{{Name: "abc"}, {Name: "abd"}},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.tree.FindMatches(tt.key, tt.limit)
			sort.Sort(ByName(actual))
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("%#v != %#v", actual, tt.expected)
			}
		})
	}
}
