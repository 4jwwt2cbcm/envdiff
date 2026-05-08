package diff

import (
	"strings"
)

// TransformFunc is a function that transforms an env value.
type TransformFunc func(value string) string

// TransformOptions controls how env values are transformed.
type TransformOptions struct {
	UppercaseKeys   bool
	LowercaseValues bool
	TrimValues      bool
	PrefixKeys      string
	StripPrefix     string
}

// TransformEnv applies transformations to an env map and returns a new map.
func TransformEnv(env map[string]string, opts TransformOptions) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		newKey := k
		newVal := v

		if opts.TrimValues {
			newVal = strings.TrimSpace(newVal)
		}
		if opts.LowercaseValues {
			newVal = strings.ToLower(newVal)
		}
		if opts.UppercaseKeys {
			newKey = strings.ToUpper(newKey)
		}
		if opts.StripPrefix != "" && strings.HasPrefix(newKey, opts.StripPrefix) {
			newKey = strings.TrimPrefix(newKey, opts.StripPrefix)
		}
		if opts.PrefixKeys != "" {
			newKey = opts.PrefixKeys + newKey
		}

		result[newKey] = newVal
	}
	return result
}

// ApplyTransformFunc applies a custom transform function to all values in an env map.
func ApplyTransformFunc(env map[string]string, fn TransformFunc) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		result[k] = fn(v)
	}
	return result
}
