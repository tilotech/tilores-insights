package stringset

// StringSet is a set of strings
type StringSet map[string]struct{}

// Add adds an element to the set and returns the set for convenience
func (s StringSet) Add(values ...string) StringSet {
	for i := range values {
		s[values[i]] = struct{}{}
	}
	return s
}

// Contains returns true if the value is inside the stringset otherwise false
func (s StringSet) Contains(value string) bool {
	_, ok := s[value]
	return ok
}

// AsSlice returns all elements as a slice
func (s StringSet) AsSlice() []string {
	keys := make([]string, len(s))
	i := 0
	for k := range s {
		keys[i] = k
		i++
	}
	return keys
}

// New creates a new StringSet
func New(values ...string) StringSet {
	set := StringSet{}
	set.Add(values...)
	return set
}
