package diff

import (
	"fmt"
	"sort"
	"strings"
)

// DependencyEdge represents a reference from one key to another.
type DependencyEdge struct {
	From string
	To   string
}

// DependencyGraph holds keys and their interpolation dependencies.
type DependencyGraph struct {
	Edges []DependencyEdge
	nodes map[string][]string // adjacency list: node -> dependencies
}

// BuildDependencyGraph inspects an env map for $VAR or ${VAR} references
// and builds a directed dependency graph.
func BuildDependencyGraph(env map[string]string) *DependencyGraph {
	g := &DependencyGraph{
		nodes: make(map[string][]string),
	}
	for key, val := range env {
		deps := extractRefs(val)
		g.nodes[key] = deps
		for _, dep := range deps {
			g.Edges = append(g.Edges, DependencyEdge{From: key, To: dep})
		}
		if _, exists := g.nodes[key]; !exists {
			g.nodes[key] = nil
		}
	}
	return g
}

// DetectCycles returns a list of cycle descriptions found in the graph.
func (g *DependencyGraph) DetectCycles() []string {
	visited := make(map[string]bool)
	inStack := make(map[string]bool)
	var cycles []string

	var dfs func(node string, path []string)
	dfs = func(node string, path []string) {
		visited[node] = true
		inStack[node] = true
		path = append(path, node)

		for _, dep := range g.nodes[node] {
			if !visited[dep] {
				dfs(dep, path)
			} else if inStack[dep] {
				// found a cycle
				start := -1
				for i, n := range path {
					if n == dep {
						start = i
						break
					}
				}
				if start >= 0 {
					cycles = append(cycles, strings.Join(path[start:], " -> ")+" -> "+dep)
				}
			}
		}
		inStack[node] = false
	}

	keys := make([]string, 0, len(g.nodes))
	for k := range g.nodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if !visited[k] {
			dfs(k, nil)
		}
	}
	return cycles
}

// TopologicalOrder returns keys in dependency-resolved order, or an error on cycles.
func (g *DependencyGraph) TopologicalOrder() ([]string, error) {
	inDegree := make(map[string]int)
	for node := range g.nodes {
		if _, ok := inDegree[node]; !ok {
			inDegree[node] = 0
		}
		for _, dep := range g.nodes[node] {
			inDegree[dep] // ensure dep exists
			inDegree[node]++ // node depends on dep, so dep must come first
			_ = inDegree[dep]
		}
	}
	// rebuild properly
	inDegree = make(map[string]int)
	for node := range g.nodes {
		if _, ok := inDegree[node]; !ok {
			inDegree[node] = 0
		}
		for _, dep := range g.nodes[node] {
			if _, ok := inDegree[dep]; !ok {
				inDegree[dep] = 0
			}
		}
	}
	for node, deps := range g.nodes {
		_ = node
		for _, dep := range deps {
			inDegree[dep]++
		}
	}

	queue := []string{}
	for node, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, node)
		}
	}
	sort.Strings(queue)

	var order []string
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		order = append(order, node)
		for _, dep := range g.nodes[node] {
			inDegree[dep]--
			if inDegree[dep] == 0 {
				queue = append(queue, dep)
				sort.Strings(queue)
			}
		}
	}
	if len(order) != len(inDegree) {
		return nil, fmt.Errorf("cycle detected: cannot produce topological order")
	}
	return order, nil
}

// extractRefs parses $VAR and ${VAR} style references from a value string.
func extractRefs(val string) []string {
	seen := map[string]bool{}
	var refs []string
	i := 0
	for i < len(val) {
		if val[i] == '$' && i+1 < len(val) {
			i++
			var key string
			if val[i] == '{' {
				i++
				start := i
				for i < len(val) && val[i] != '}' {
					i++
				}
				key = val[start:i]
				if i < len(val) {
					i++
				}
			} else {
				start := i
				for i < len(val) && (isAlnum(val[i]) || val[i] == '_') {
					i++
				}
				key = val[start:i]
			}
			if key != "" && !seen[key] {
				seen[key] = true
				refs = append(refs, key)
			}
		} else {
			i++
		}
	}
	return refs
}

func isAlnum(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}
