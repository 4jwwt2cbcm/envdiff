package diff

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// PatchOp represents a single patch operation.
type PatchOp struct {
	Key    string `json:"key"`
	Action string `json:"action"` // "add", "remove", "update"
	OldVal string `json:"old_value,omitempty"`
	NewVal string `json:"new_value,omitempty"`
}

// Patch represents a set of operations to transform one env into another.
type Patch struct {
	Ops []PatchOp `json:"ops"`
}

// GeneratePatch produces a Patch that, when applied to src, yields dst.
func GeneratePatch(src, dst map[string]string) Patch {
	var ops []PatchOp

	for k, dv := range dst {
		if sv, ok := src[k]; !ok {
			ops = append(ops, PatchOp{Key: k, Action: "add", NewVal: dv})
		} else if sv != dv {
			ops = append(ops, PatchOp{Key: k, Action: "update", OldVal: sv, NewVal: dv})
		}
	}

	for k, sv := range src {
		if _, ok := dst[k]; !ok {
			ops = append(ops, PatchOp{Key: k, Action: "remove", OldVal: sv})
		}
	}

	sort.Slice(ops, func(i, j int) bool {
		return ops[i].Key < ops[j].Key
	})

	return Patch{Ops: ops}
}

// ApplyPatch applies a Patch to env, returning the modified copy.
func ApplyPatch(env map[string]string, p Patch) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	for _, op := range p.Ops {
		switch op.Action {
		case "add":
			out[op.Key] = op.NewVal
		case "remove":
			delete(out, op.Key)
		case "update":
			out[op.Key] = op.NewVal
		default:
			return nil, fmt.Errorf("unknown patch action: %q", op.Action)
		}
	}

	return out, nil
}

// WritePatch writes a patch as a human-readable diff to the given file path.
func WritePatch(path string, p Patch) error {
	var sb strings.Builder
	for _, op := range p.Ops {
		switch op.Action {
		case "add":
			fmt.Fprintf(&sb, "+ %s=%s\n", op.Key, op.NewVal)
		case "remove":
			fmt.Fprintf(&sb, "- %s=%s\n", op.Key, op.OldVal)
		case "update":
			fmt.Fprintf(&sb, "~ %s: %s -> %s\n", op.Key, op.OldVal, op.NewVal)
		}
	}
	return os.WriteFile(path, []byte(sb.String()), 0644)
}
