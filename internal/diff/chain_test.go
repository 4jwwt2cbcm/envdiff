package diff

import (
	"errors"
	"strings"
	"testing"
)

func baseChainEnv() map[string]string {
	return map[string]string{
		"APP_HOST": "  localhost  ",
		"app_port": "8080",
		"DB_PASS":  "secret",
	}
}

func TestChain_SingleStep(t *testing.T) {
	env := baseChainEnv()
	steps := []ChainStep{
		{Name: "trim", Fn: func(m map[string]string) (map[string]string, error) {
			out := map[string]string{}
			for k, v := range m {
				out[k] = strings.TrimSpace(v)
			}
			return out, nil
		}},
	}
	results, err := Chain(env, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Output["APP_HOST"] != "localhost" {
		t.Errorf("expected trimmed value, got %q", results[0].Output["APP_HOST"])
	}
}

func TestChain_MultipleSteps(t *testing.T) {
	env := baseChainEnv()
	steps := []ChainStep{
		{Name: "uppercase-keys", Fn: func(m map[string]string) (map[string]string, error) {
			out := map[string]string{}
			for k, v := range m {
				out[strings.ToUpper(k)] = v
			}
			return out, nil
		}},
		{Name: "add-prefix", Fn: func(m map[string]string) (map[string]string, error) {
			out := map[string]string{}
			for k, v := range m {
				out["ENV_"+k] = v
			}
			return out, nil
		}},
	}
	results, err := Chain(env, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := results[1].Output["ENV_APP_PORT"]; !ok {
		t.Errorf("expected ENV_APP_PORT in final output")
	}
}

func TestChain_StopsOnError(t *testing.T) {
	env := baseChainEnv()
	ran := false
	steps := []ChainStep{
		{Name: "fail", Fn: func(m map[string]string) (map[string]string, error) {
			return nil, errors.New("step failed")
		}},
		{Name: "never", Fn: func(m map[string]string) (map[string]string, error) {
			ran = true
			return m, nil
		}},
	}
	_, err := Chain(env, steps)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if ran {
		t.Error("second step should not have run after failure")
	}
}

func TestChainSummary_KeyChanges(t *testing.T) {
	original := map[string]string{"A": "1", "B": "2"}
	results := []ChainResult{
		{Step: "add-C", Output: map[string]string{"A": "1", "B": "2", "C": "3"}},
		{Step: "remove-A", Output: map[string]string{"B": "2", "C": "3"}},
	}
	lines := ChainSummary(results, original)
	if len(lines) != 2 {
		t.Fatalf("expected 2 summary lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "+1") {
		t.Errorf("expected +1 key in step 1, got: %s", lines[0])
	}
	if !strings.Contains(lines[1], "-1") {
		t.Errorf("expected -1 key in step 2, got: %s", lines[1])
	}
}
