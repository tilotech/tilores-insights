package stringset

// Compare creates a diff between the stringset s1 and s2 and returns only the
// entries from s1 that are not present in s2.
func Compare(s1, s2 []string) []string {
	m := make(map[string]int)

	for _, v2 := range s2 {
		m[v2]++
	}

	ret := []string{}
	for _, v1 := range s1 {
		if m[v1] > 0 {
			m[v1]--
			continue
		}
		ret = append(ret, v1)
	}

	return ret
}
