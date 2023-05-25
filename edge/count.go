package edge

import api "github.com/tilotech/tilores-plugin-api"

// Count returns the amount of edges in the provided list.
//
// This function does not consider implicit edges based on duplicates.
func Count(edges api.Edges) int {
	return len(edges)
}
