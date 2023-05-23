package edge

import (
	"fmt"
	"strings"

	api "github.com/tilotech/tilores-plugin-api"
)

// MatrixEntry represents a single entry for an edge and duplicate matrix.
type MatrixEntry struct {
	A     string          `json:"a"`
	B     string          `json:"b"`
	Links map[string]bool `json:"links"`
}

// Matrix returns a matrix in which it is possible to see the links between each
// two records and due to which rule or duplicate they are linked.
func Matrix(edges api.Edges, duplicates api.Duplicates, links []string) []MatrixEntry {
	matrix := map[string]map[string]bool{}

	if links == nil {
		links = collectLinks(edges, duplicates)
	}

	validLinks := newSet(links...)

	for _, edge := range edges {
		a, b, rule := splitEdge(edge)
		if !validLinks.Contains(rule) {
			continue
		}
		k := fmt.Sprintf("%s:%s", a, b)
		if _, ok := matrix[k]; !ok {
			matrix[k] = initializeMatrixEntry(links)
		}
		matrix[k][rule] = true
	}

	for id1, dups := range duplicates {
		rule := "duplicate"
		parts := strings.SplitN(id1, ":", 2)
		if len(parts) == 2 {
			rule = fmt.Sprintf("%s:%s", parts[0], rule)
			id1 = parts[1]
		}
		if !validLinks.Contains(rule) {
			continue
		}
		for _, id2 := range dups {
			a, b := id1, id2
			if b < a {
				a, b = b, a
			}
			k := fmt.Sprintf("%s:%s", a, b)
			if _, ok := matrix[k]; !ok {
				matrix[k] = initializeMatrixEntry(links)
			}
			matrix[k][rule] = true
		}
	}

	entries := make([]MatrixEntry, 0, len(matrix))
	for k, l := range matrix {
		parts := strings.SplitN(k, ":", 2)
		entries = append(entries, MatrixEntry{
			A:     parts[0],
			B:     parts[1],
			Links: l,
		})
	}

	return entries
}

func splitEdge(edge string) (string, string, string) {
	parts := strings.SplitN(edge, ":", 3)
	if parts[0] < parts[1] {
		return parts[0], parts[1], parts[2]
	}
	return parts[1], parts[0], parts[2]
}

func initializeMatrixEntry(links []string) map[string]bool {
	entry := make(map[string]bool, len(links))
	for _, l := range links {
		entry[l] = false
	}
	return entry
}

func collectLinks(edges api.Edges, duplicates api.Duplicates) []string {
	s := newSet[string]()
	for _, edge := range edges {
		_, _, rule := splitEdge(edge)
		s.Add(rule)
	}
	for id := range duplicates {
		rule := "duplicate"
		parts := strings.SplitN(id, ":", 2)
		if len(parts) == 2 {
			rule = fmt.Sprintf("%s:%s", parts[0], rule)
		}
		s.Add(rule)
	}
	return s.AsSlice()
}
