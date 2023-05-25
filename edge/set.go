package edge

type set[T comparable] map[T]struct{}

func (s set[T]) Add(values ...T) set[T] {
	for i := range values {
		s[values[i]] = struct{}{}
	}
	return s
}

func (s set[T]) Contains(value T) bool {
	_, ok := s[value]
	return ok
}

func (s set[T]) AsSlice() []T {
	keys := make([]T, len(s))
	i := 0
	for k := range s {
		keys[i] = k
		i++
	}
	return keys
}

func newSet[T comparable](values ...T) set[T] {
	s := make(set[T], len(values))
	return s.Add(values...)
}
