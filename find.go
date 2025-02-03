package raildata

import (
	"strings"
)

// Finder is an interface to find objects of type T with codes of type C.
type Finder[T any, C ~string] interface {
	// Sets up to find an object with the given code.
	WithCode(code C) Finder[T, C]
	// Sets up to find an object with the given name.
	WithName(name string) Finder[T, C]
	// Searches for the object, returning either (non-nil pointer, true) if found,
	// or (nil pointer, false) if not found.
	Search() (*T, bool)
	// Searches for the object, returning either the found object, or a made-up
	// object that was built from the search data.
	SearchOrSynthesize() *T
}

type finderImpl[T any, C ~string] struct {
	byCode        map[string]*T
	byName        map[string]*T
	byAbbr        map[string]*T
	list          []T
	getCandidates func(s *T) []string
	synthesize    func(code *C, name *string) *T
	code          *C
	name          *string
}

func (f finderImpl[T, C]) WithCode(code C) Finder[T, C] {
	f.code = &code
	return f
}

func (f finderImpl[T, C]) WithName(name string) Finder[T, C] {
	f.name = &name
	return f
}

func (f finderImpl[T, C]) Search() (*T, bool) {
	if f.code != nil {
		codeLc := strings.ToLower(string(*f.code))
		if item, found := f.byCode[codeLc]; found {
			return item, true
		}
	}
	if f.name != nil {
		nameLc := strings.ToLower(*f.name)
		if item, found := f.byName[nameLc]; found {
			return item, true
		}
		if item, found := f.byAbbr[nameLc]; found {
			return item, true
		}
		item, matchLen := fuzzyFind(nameLc, f.list, f.getCandidates)
		if matchLen > 2 && matchLen >= len(nameLc)/4 {
			return item, true
		}
	}
	return nil, false
}

func (f finderImpl[T, C]) SearchOrSynthesize() *T {
	if item, found := f.Search(); found {
		return item
	}
	return f.synthesize(f.code, f.name)
}

func fuzzyFind[T any](input string, list []T, getCandidates func(*T) []string) (best *T, matchLen int) {
	best = nil
	matchLen = 0
	strLen := 0
	for i := range list {
		for _, candidate := range getCandidates(&list[i]) {
			candidate := strings.ToLower(candidate)
			ml := fuzzyMatch(input, candidate)
			if ml > matchLen || (ml == matchLen && len(candidate) < strLen) {
				best = &list[i]
				matchLen = ml
				strLen = len(candidate)
			}
		}
	}
	return
}

func fuzzyMatch(input string, candidate string) int {
	rs := len(input) + 1
	cs := len(candidate) + 1
	matchLen := make([]int, rs*cs)
	for row := 1; row < cs; row++ {
		for col := 1; col < rs; col++ {
			idx := row*rs + col
			m := 0
			if input[col-1] == candidate[row-1] {
				m = 1 + matchLen[idx-rs-1]
			}
			m = max(m, matchLen[idx-rs], matchLen[idx-1], matchLen[idx-rs-1])
			matchLen[idx] = m
		}
	}
	return matchLen[rs*cs-1]
}
