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

func makeLeaf(value Location) *Trie {
	tree := NewTrie()
	tree.leaf = true
	tree.value = value
	return tree
}

func TestTrie_Insert(t *testing.T) {
	tests := map[string]struct {
		before *Trie
		key    string
		value  Location
		after  *Trie
	}{
		"single character empty tree": {
			NewTrie(),
			"a",
			Location{Name: "a"},
			makeTree('a', makeLeaf(Location{Name: "a"})),
		},
		"empty key empty tree": {
			NewTrie(),
			"",
			Location{Name: ""},
			makeLeaf(Location{Name: ""}),
		},
		"multiple characters empty tree": {
			NewTrie(),
			"abc",
			Location{Name: "abc"},
			makeTree('a', makeTree('b', makeTree('c', makeLeaf(Location{Name: "abc"})))),
		},
		"multiple characters case-insensitive": {
			NewTrie(),
			"ABC",
			Location{Name: "ABC"},
			makeTree('a', makeTree('b', makeTree('c', makeLeaf(Location{Name: "ABC"})))),
		},
		"multiple characters non-empty tree": {
			makeTree('a', makeTree('b', makeTree('c', makeLeaf(Location{Name: "abc"})))),
			"abd",
			Location{Name: "abd"},
			makeTree(
				'a', makeTree(
					'b', makeTree(
						'c', makeLeaf(Location{Name: "abc"}),
						'd', makeLeaf(Location{Name: "abd"}),
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
			makeTree('a', makeLeaf(Location{Name: "a"})),
			"nope",
			false,
		},
		"single character": {
			makeTree('a', makeLeaf(Location{Name: "a"})),
			"a",
			true,
		},
		"multiple characters": {
			makeTree('a', makeTree('b', makeTree('c', makeLeaf(Location{Name: "abc"})))),
			"abc",
			true,
		},
		"multiple character subset": {
			makeTree('a', makeTree('b', makeTree('c', makeLeaf(Location{Name: "abc"})))),
			"ab",
			false,
		},
		"multiple character superset": {
			makeTree('a', makeTree('b', makeTree('c', makeLeaf(Location{Name: "abc"})))),
			"abcd",
			false,
		},
		"case-insensitive": {
			makeTree('a', makeTree('b', makeTree('c', makeLeaf(Location{Name: "abc"})))),
			"ABC",
			true,
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

func TestTrie_FindMatches(t *testing.T) {
	tests := map[string]struct {
		tree     *Trie
		key      string
		limit    int
		expected []Location
	}{
		"empty tree": {
			NewTrie(),
			"nope",
			10,
			[]Location{},
		},
		"missing key": {
			makeTree('a', makeLeaf(Location{Name: "a"})),
			"nope",
			10,
			[]Location{},
		},
		"exact match": {
			makeTree('a', makeLeaf(Location{Name: "a"})),
			"a",
			10,
			[]Location{{Name: "a"}},
		},
		"multiple matches": {
			makeTree(
				'a', makeTree(
					'b', makeTree(
						'c', makeLeaf(Location{Name: "abc"}),
						'd', makeLeaf(Location{Name: "abd"}),
					),
				),
			),
			"ab",
			10,
			[]Location{{Name: "abc"}, {Name: "abd"}},
		},
		"case-insensitive": {
			makeTree(
				'a', makeTree(
					'b', makeTree(
						'c', makeLeaf(Location{Name: "ABC"}),
					),
				),
			),
			"ABC",
			10,
			[]Location{{Name: "ABC"}},
		},
		"multiple matches limit returns shortest first": {
			makeTree('a', makeTree('b', makeTree(
				'c', makeLeaf(Location{Name: "abc"}),
				'd', makeTree('e', makeLeaf(Location{Name: "abde"})),
			))),
			"ab",
			1,
			[]Location{{Name: "abc"}},
		},
		"limit < 0 means no limit": {
			makeTree(
				'a', makeTree(
					'b', makeTree(
						'c', makeLeaf(Location{Name: "abc"}),
						'd', makeLeaf(Location{Name: "abd"}),
					),
				),
			),
			"ab",
			-1,
			[]Location{{Name: "abc"}, {Name: "abd"}},
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
