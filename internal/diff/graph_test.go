package diff

import (
	"strings"
	"testing"
)

func TestBuildDependencyGraph_NoRefs(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
	}
	g := BuildDependencyGraph(env)
	if len(g.Edges) != 0 {
		t.Errorf("expected 0 edges, got %d", len(g.Edges))
	}
}

func TestBuildDependencyGraph_BraceRef(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "http://${HOST}:${PORT}",
		"HOST":     "localhost",
		"PORT":     "8080",
	}
	g := BuildDependencyGraph(env)
	if len(g.Edges) != 2 {
		t.Errorf("expected 2 edges, got %d", len(g.Edges))
	}
	for _, e := range g.Edges {
		if e.From != "BASE_URL" {
			t.Errorf("expected edge from BASE_URL, got %s", e.From)
		}
	}
}

func TestBuildDependencyGraph_BareRef(t *testing.T) {
	env := map[string]string{
		"FULL_URL": "$SCHEME://$HOST",
		"SCHEME":   "https",
		"HOST":     "example.com",
	}
	g := BuildDependencyGraph(env)
	if len(g.Edges) != 2 {
		t.Errorf("expected 2 edges, got %d", len(g.Edges))
	}
}

func TestDetectCycles_NoCycle(t *testing.T) {
	env := map[string]string{
		"A": "$B",
		"B": "$C",
		"C": "literal",
	}
	g := BuildDependencyGraph(env)
	cycles := g.DetectCycles()
	if len(cycles) != 0 {
		t.Errorf("expected no cycles, got: %v", cycles)
	}
}

func TestDetectCycles_WithCycle(t *testing.T) {
	env := map[string]string{
		"A": "$B",
		"B": "$A",
	}
	g := BuildDependencyGraph(env)
	cycles := g.DetectCycles()
	if len(cycles) == 0 {
		t.Error("expected a cycle to be detected")
	}
	found := false
	for _, c := range cycles {
		if strings.Contains(c, "A") && strings.Contains(c, "B") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected cycle involving A and B, got: %v", cycles)
	}
}

func TestTopologicalOrder_Simple(t *testing.T) {
	env := map[string]string{
		"URL":  "http://$HOST:$PORT",
		"HOST": "localhost",
		"PORT": "8080",
	}
	g := BuildDependencyGraph(env)
	order, err := g.TopologicalOrder()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// HOST and PORT must appear before URL
	pos := map[string]int{}
	for i, k := range order {
		pos[k] = i
	}
	if pos["URL"] < pos["HOST"] || pos["URL"] < pos["PORT"] {
		t.Errorf("URL should come after HOST and PORT in order: %v", order)
	}
}

func TestTopologicalOrder_CycleReturnsError(t *testing.T) {
	env := map[string]string{
		"X": "$Y",
		"Y": "$X",
	}
	g := BuildDependencyGraph(env)
	_, err := g.TopologicalOrder()
	if err == nil {
		t.Error("expected error for cyclic graph")
	}
}

func TestExtractRefs_Mixed(t *testing.T) {
	refs := extractRefs("${FOO} and $BAR plus ${BAZ}")
	if len(refs) != 3 {
		t.Errorf("expected 3 refs, got %d: %v", len(refs), refs)
	}
}

func TestExtractRefs_NoDuplicates(t *testing.T) {
	refs := extractRefs("$A $A ${A}")
	if len(refs) != 1 {
		t.Errorf("expected 1 unique ref, got %d: %v", len(refs), refs)
	}
}
