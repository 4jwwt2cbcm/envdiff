package diff

import (
	"strings"
	"testing"
)

func baseTransformEnv() map[string]string {
	return map[string]string{
		"APP_HOST": "  localhost  ",
		"APP_PORT": "8080",
		"DB_URL":   "Postgres://localhost",
	}
}

func TestTransformEnv_TrimValues(t *testing.T) {
	result := TransformEnv(baseTransformEnv(), TransformOptions{TrimValues: true})
	if result["APP_HOST"] != "localhost" {
		t.Errorf("expected trimmed value, got %q", result["APP_HOST"])
	}
}

func TestTransformEnv_LowercaseValues(t *testing.T) {
	result := TransformEnv(baseTransformEnv(), TransformOptions{LowercaseValues: true})
	if result["DB_URL"] != "postgres://localhost" {
		t.Errorf("expected lowercase value, got %q", result["DB_URL"])
	}
}

func TestTransformEnv_UppercaseKeys(t *testing.T) {
	env := map[string]string{"app_host": "localhost", "app_port": "8080"}
	result := TransformEnv(env, TransformOptions{UppercaseKeys: true})
	if _, ok := result["APP_HOST"]; !ok {
		t.Error("expected APP_HOST key after uppercase transform")
	}
	if _, ok := result["app_host"]; ok {
		t.Error("expected original lowercase key to be gone")
	}
}

func TestTransformEnv_PrefixKeys(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	result := TransformEnv(env, TransformOptions{PrefixKeys: "APP_"})
	if _, ok := result["APP_HOST"]; !ok {
		t.Error("expected APP_HOST after prefix transform")
	}
	if _, ok := result["APP_PORT"]; !ok {
		t.Error("expected APP_PORT after prefix transform")
	}
}

func TestTransformEnv_StripPrefix(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "OTHER": "x"}
	result := TransformEnv(env, TransformOptions{StripPrefix: "APP_"})
	if _, ok := result["HOST"]; !ok {
		t.Error("expected HOST after strip prefix")
	}
	if _, ok := result["OTHER"]; !ok {
		t.Error("expected OTHER (no prefix) to remain")
	}
}

func TestApplyTransformFunc_UppercaseValues(t *testing.T) {
	env := map[string]string{"KEY": "hello", "OTHER": "world"}
	result := ApplyTransformFunc(env, strings.ToUpper)
	if result["KEY"] != "HELLO" {
		t.Errorf("expected HELLO, got %q", result["KEY"])
	}
	if result["OTHER"] != "WORLD" {
		t.Errorf("expected WORLD, got %q", result["OTHER"])
	}
}

func TestTransformEnv_NoOptions(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	result := TransformEnv(env, TransformOptions{})
	if result["KEY"] != "value" {
		t.Errorf("expected unchanged value, got %q", result["KEY"])
	}
}
