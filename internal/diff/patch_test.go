package diff

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGeneratePatch_AddRemoveUpdate(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{"A": "1", "B": "99", "D": "4"}

	p := GeneratePatch(src, dst)

	actions := map[string]PatchOp{}
	for _, op := range p.Ops {
		actions[op.Key] = op
	}

	if op, ok := actions["B"]; !ok || op.Action != "update" || op.OldVal != "2" || op.NewVal != "99" {
		t.Errorf("expected update for B, got %+v", op)
	}
	if op, ok := actions["C"]; !ok || op.Action != "remove" {
		t.Errorf("expected remove for C, got %+v", op)
	}
	if op, ok := actions["D"]; !ok || op.Action != "add" || op.NewVal != "4" {
		t.Errorf("expected add for D, got %+v", op)
	}
	if _, ok := actions["A"]; ok {
		t.Error("A should not appear in patch (unchanged)")
	}
}

func TestGeneratePatch_NoDiff(t *testing.T) {
	env := map[string]string{"X": "1", "Y": "2"}
	p := GeneratePatch(env, env)
	if len(p.Ops) != 0 {
		t.Errorf("expected empty patch, got %d ops", len(p.Ops))
	}
}

func TestApplyPatch_Basic(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"A": "1", "B": "99", "C": "3"}

	p := GeneratePatch(src, dst)
	result, err := ApplyPatch(src, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range dst {
		if result[k] != v {
			t.Errorf("key %s: want %q, got %q", k, v, result[k])
		}
	}
}

func TestApplyPatch_UnknownAction(t *testing.T) {
	p := Patch{Ops: []PatchOp{{Key: "X", Action: "invalid"}}}
	_, err := ApplyPatch(map[string]string{}, p)
	if err == nil {
		t.Error("expected error for unknown action")
	}
}

func TestWritePatch(t *testing.T) {
	p := Patch{Ops: []PatchOp{
		{Key: "A", Action: "add", NewVal: "hello"},
		{Key: "B", Action: "remove", OldVal: "bye"},
		{Key: "C", Action: "update", OldVal: "x", NewVal: "y"},
	}}

	tmp := filepath.Join(t.TempDir(), "patch.txt")
	if err := WritePatch(tmp, p); err != nil {
		t.Fatalf("WritePatch error: %v", err)
	}

	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}

	content := string(data)
	for _, want := range []string{"+ A=hello", "- B=bye", "~ C: x -> y"} {
		if !containsStr(content, want) {
			t.Errorf("patch output missing %q\ngot:\n%s", want, content)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
