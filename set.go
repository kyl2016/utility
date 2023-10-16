package utility

type IntSet struct {
	m map[int]bool
}

func NewIntSet() IntSet {
	s := IntSet{}
	s.m = make(map[int]bool)
	return s
}

func (s *IntSet) Add(values ...int) {
	for _, value := range values {
		s.m[value] = true
	}
}

func (s *IntSet) Remove(values ...int) {
	for _, value := range values {
		delete(s.m, value)
	}
}

func (s *IntSet) Contains(value int) bool {
	_, c := s.m[value]
	return c
}

func (s *IntSet) All() []int {
	all := make([]int, len(s.m))
	i := 0
	for key := range s.m {
		all[i] = key
		i++
	}
	return all
}

func (s *IntSet) Count() int {
	return len(s.m)
}

type StringSet struct {
	m map[string]bool
}

func NewStringSet() StringSet {
	s := StringSet{}
	s.m = make(map[string]bool)
	return s
}

func (s *StringSet) Add(values ...string) {
	for _, value := range values {
		s.m[value] = true
	}
}

func (s *StringSet) Remove(values ...string) {
	for _, value := range values {
		delete(s.m, value)
	}
}

func (s *StringSet) Contains(value string) bool {
	_, c := s.m[value]
	return c
}

func (s *StringSet) All() []string {
	all := make([]string, len(s.m))
	i := 0
	for key := range s.m {
		all[i] = key
		i++
	}
	return all
}

func (s *StringSet) Count() int {
	return len(s.m)
}
