package fsm

import (
	"testing"
)

func TestSet(t *testing.T) {
	var s states = make([]string, 0)
	s.add("def")
	s.add("abc")
	s.add("def")
	if s[0] != "abc" || s[1] != "def" {
		t.Fatal("Expected just 'abc' and 'def' in states")
	}
	s.remove("def")
	s.remove("xyz")
	if s[0] != "abc" || len(s) != 1 {
		t.Fatal("Expected just 'abc' in states")
	}
}
