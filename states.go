package fsm

import (
	"sort"
	"strings"
)

type states []string

func (ss *states) Len() int {
	return len(*ss)
}

func (ss *states) Less(i, j int) bool {
	return strings.Compare((*ss)[i], (*ss)[j]) < 0
}

func (ss *states) Swap(i, j int) {
	(*ss)[i], (*ss)[j] = (*ss)[j], (*ss)[i]
}

func (ss *states) add(s string) {
	for i := range *ss {
		if (*ss)[i] == s {
			return
		}
	}
	*ss = append(*ss, s)
	sort.Sort(ss)
}

func (ss *states) addmany(s []string) {
	if s == nil {
		return
	}
	for i := range s {
		ss.add(s[i])
	}
	sort.Sort(ss)
}

func (ss *states) remove(s string) {
	for i := range *ss {
		if (*ss)[i] == s {
			*ss = append((*ss)[:i], (*ss)[i+1:]...)
		}
	}
	sort.Sort(ss)
}

func (ss *states) equal(ss2 states) bool {
	if len(*ss) != len(ss2) {
		return false
	}
	for i := range *ss {
		if (*ss)[i] != ss2[i] {
			return false
		}
	}
	return true
}
