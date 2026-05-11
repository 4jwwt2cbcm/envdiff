package diff

import (
	"testing"
)

func baseNormalizeEnv() map[string]string {
	return map[string]string{
		"APP_NAME": "  myapp  ",
		"DB_HOST":  "localhost",
		"EMPTY_KEY": "",
		"api_key":  "secret",
	}
}

func TestNormalizeEnv_NoOptions(t *testing.T) {
	env := baseNormalizeEnv()
	res := NormalizeEnv(env, NormalizeOptions{})

	if len(res.Changed) != 0 {
		t.Errorf("expected no changes, got %v", res.Changed)
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected no removals, got %v", res.Removed)
	}
	if res.Env["APP_NAME"] != "  myapp  " {
		t.Errorf("expected untrimmed value, got %q", res.Env["APP_NAME"])
	}
}

func TestNormalizeEnv_TrimValues(t *testing.T) {
	env := baseNormalizeEnv()
	res := NormalizeEnv(env, NormalizeOptions{TrimValues: true})

	if res.Env["APP_NAME"] != "myapp" {
		t.Errorf("expected trimmed value, got %q", res.Env["APP_NAME"])
	}
	if res.Env["DB_HOST"] != "localhost" {
		t.Errorf("expected unchanged value, got %q", res.Env["DB_HOST"])
	}
	if len(res.Changed) != 1 || res.Changed[0] != "APP_NAME" {
		t.Errorf("expected APP_NAME in changed, got %v", res.Changed)
	}
}

func TestNormalizeEnv_UppercaseKeys(t *testing.T) {
	env := baseNormalizeEnv()
	res := NormalizeEnv(env, NormalizeOptions{UppercaseKeys: true})

	if _, ok := res.Env["API_KEY"]; !ok {
		t.Errorf("expected API_KEY to exist after uppercase normalization")
	}
	if _, ok := res.Env["api_key"]; ok {
		t.Errorf("expected api_key to be renamed")
	}
}

func TestNormalizeEnv_LowercaseKeys(t *testing.T) {
	env := baseNormalizeEnv()
	res := NormalizeEnv(env, NormalizeOptions{LowercaseKeys: true})

	if _, ok := res.Env["app_name"]; !ok {
		t.Errorf("expected app_name to exist after lowercase normalization")
	}
}

func TestNormalizeEnv_RemoveEmpty(t *testing.T) {
	env := baseNormalizeEnv()
	res := NormalizeEnv(env, NormalizeOptions{RemoveEmpty: true})

	if _, ok := res.Env["EMPTY_KEY"]; ok {
		t.Errorf("expected EMPTY_KEY to be removed")
	}
	if len(res.Removed) != 1 || res.Removed[0] != "EMPTY_KEY" {
		t.Errorf("expected EMPTY_KEY in removed, got %v", res.Removed)
	}
}

func TestNormalizeEnv_TrimAndRemoveEmpty(t *testing.T) {
	env := map[string]string{
		"KEY_A": "  ",
		"KEY_B": " value ",
	}
	res := NormalizeEnv(env, NormalizeOptions{TrimValues: true, RemoveEmpty: true})

	if _, ok := res.Env["KEY_A"]; ok {
		t.Errorf("expected KEY_A to be removed after trim+removeEmpty")
	}
	if res.Env["KEY_B"] != "value" {
		t.Errorf("expected trimmed KEY_B, got %q", res.Env["KEY_B"])
	}
}
