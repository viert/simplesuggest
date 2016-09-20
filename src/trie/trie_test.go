package trie

import (
	"github.com/nu7hatch/gouuid"
	"reflect"
	"sort"
	"testing"
)

func TestNodeConsistency(t *testing.T) {

	var really []string
	var expected []string
	var prefix string

	root := newNode()
	root.insert("abc", "12", 0)
	root.insert("abcd", "15", 0)
	root.insert("abr-", "18", 0)

	really = root.getAll()
	expected = []string{"12", "15", "18"}
	sort.Strings(really)
	if !reflect.DeepEqual(really, expected) {
		t.Errorf("Invalid tree getAll() after insertion: %v instead of %v", really, expected)
	}

	prefix = "a"
	really = root.getAllWithPrefix(prefix, 0)
	expected = []string{"12", "15", "18"}
	sort.Strings(really)
	if !reflect.DeepEqual(really, expected) {
		t.Errorf("Invalid tree getAllWithPrefix('%s'): %v instead of %v", prefix, really, expected)
	}

	prefix = "ab"
	really = root.getAllWithPrefix(prefix, 0)
	expected = []string{"12", "15", "18"}
	sort.Strings(really)
	if !reflect.DeepEqual(really, expected) {
		t.Errorf("Invalid tree getAllWithPrefix('%s'): %v instead of %v", prefix, really, expected)
	}

	prefix = "abc"
	really = root.getAllWithPrefix(prefix, 0)
	expected = []string{"12", "15"}
	sort.Strings(really)
	if !reflect.DeepEqual(really, expected) {
		t.Errorf("Invalid tree getAllWithPrefix('%s'): %v instead of %v", prefix, really, expected)
	}

	prefix = "abcd"
	really = root.getAllWithPrefix(prefix, 0)
	expected = []string{"15"}
	sort.Strings(really)
	if !reflect.DeepEqual(really, expected) {
		t.Errorf("Invalid tree getAllWithPrefix('%s'): %v instead of %v", prefix, really, expected)
	}

	prefix = "abr-"
	really = root.getAllWithPrefix(prefix, 0)
	expected = []string{"18"}
	sort.Strings(really)
	if !reflect.DeepEqual(really, expected) {
		t.Errorf("Invalid tree getAllWithPrefix('%s'): %v instead of %v", prefix, really, expected)
	}

	prefix = "blahminor"
	really = root.getAllWithPrefix(prefix, 0)
	expected = []string{}
	sort.Strings(really)
	if !reflect.DeepEqual(really, expected) {
		t.Errorf("Invalid tree getAllWithPrefix('%s'): %v instead of %v", prefix, really, expected)
	}

}

func BenchmarkNodeInsert(b *testing.B) {
	t := NewTrie()
	value := "obj"
	for i := 0; i < b.N; i++ {
		u, _ := uuid.NewV4()
		key := u.String()
		t.Insert(key, value)
	}
}
