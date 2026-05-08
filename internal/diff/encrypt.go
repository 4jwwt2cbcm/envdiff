package diff

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// EncryptOptions controls encryption behavior.
type EncryptOptions struct {
	Passphrase string
}

// EncryptedEnv holds an encrypted representation of an env map.
type EncryptedEnv struct {
	Data map[string]string `json:"data"`
}

// deriveKey produces a 32-byte AES key from a passphrase via SHA-256.
func deriveKey(passphrase string) []byte {
	h := sha256.Sum256([]byte(passphrase))
	return h[:]
}

// EncryptEnv encrypts all values in the env map using AES-GCM.
// Each value is independently encrypted and base64-encoded.
func EncryptEnv(env map[string]string, opts EncryptOptions) (EncryptedEnv, error) {
	if opts.Passphrase == "" {
		return EncryptedEnv{}, errors.New("encrypt: passphrase must not be empty")
	}
	key := deriveKey(opts.Passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return EncryptedEnv{}, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return EncryptedEnv{}, err
	}

	result := EncryptedEnv{Data: make(map[string]string, len(env))}
	for k, v := range env {
		nonce := make([]byte, gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return EncryptedEnv{}, err
		}
		ciphertext := gcm.Seal(nonce, nonce, []byte(v), nil)
		result.Data[k] = base64.StdEncoding.EncodeToString(ciphertext)
	}
	return result, nil
}

// DecryptEnv decrypts an EncryptedEnv back to a plain env map.
func DecryptEnv(enc EncryptedEnv, opts EncryptOptions) (map[string]string, error) {
	if opts.Passphrase == "" {
		return nil, errors.New("decrypt: passphrase must not be empty")
	}
	key := deriveKey(opts.Passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(enc.Data))
	for k, v := range enc.Data {
		ciphertext, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return nil, err
		}
		if len(ciphertext) < gcm.NonceSize() {
			return nil, errors.New("decrypt: ciphertext too short for key " + k)
		}
		nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return nil, errors.New("decrypt: failed to decrypt key " + k + ": " + err.Error())
		}
		result[k] = string(plaintext)
	}
	return result, nil
}
