package diff

import (
	"testing"
)

func sampleEnvForEncrypt() map[string]string {
	return map[string]string{
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "abc123",
		"APP_ENV":     "production",
	}
}

func TestEncryptEnv_RoundTrip(t *testing.T) {
	env := sampleEnvForEncrypt()
	opts := EncryptOptions{Passphrase: "test-passphrase"}

	enc, err := EncryptEnv(env, opts)
	if err != nil {
		t.Fatalf("EncryptEnv error: %v", err)
	}

	if len(enc.Data) != len(env) {
		t.Errorf("expected %d keys, got %d", len(env), len(enc.Data))
	}

	dec, err := DecryptEnv(enc, opts)
	if err != nil {
		t.Fatalf("DecryptEnv error: %v", err)
	}

	for k, want := range env {
		got, ok := dec[k]
		if !ok {
			t.Errorf("missing key %q after decrypt", k)
			continue
		}
		if got != want {
			t.Errorf("key %q: got %q, want %q", k, got, want)
		}
	}
}

func TestEncryptEnv_ValuesAreObfuscated(t *testing.T) {
	env := map[string]string{"SECRET": "plaintext"}
	opts := EncryptOptions{Passphrase: "mypass"}

	enc, err := EncryptEnv(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enc.Data["SECRET"] == "plaintext" {
		t.Error("encrypted value should not equal plaintext")
	}
}

func TestEncryptEnv_EmptyPassphrase(t *testing.T) {
	_, err := EncryptEnv(map[string]string{"K": "V"}, EncryptOptions{})
	if err == nil {
		t.Error("expected error for empty passphrase")
	}
}

func TestDecryptEnv_WrongPassphrase(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	enc, err := EncryptEnv(env, EncryptOptions{Passphrase: "correct"})
	if err != nil {
		t.Fatalf("encrypt error: %v", err)
	}
	_, err = DecryptEnv(enc, EncryptOptions{Passphrase: "wrong"})
	if err == nil {
		t.Error("expected error when decrypting with wrong passphrase")
	}
}

func TestDecryptEnv_EmptyPassphrase(t *testing.T) {
	_, err := DecryptEnv(EncryptedEnv{Data: map[string]string{"K": "V"}}, EncryptOptions{})
	if err == nil {
		t.Error("expected error for empty passphrase")
	}
}

func TestEncryptEnv_EmptyMap(t *testing.T) {
	enc, err := EncryptEnv(map[string]string{}, EncryptOptions{Passphrase: "pass"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(enc.Data) != 0 {
		t.Errorf("expected empty encrypted map")
	}
}
