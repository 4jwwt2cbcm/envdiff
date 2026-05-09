package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// Profile represents a named set of keys to focus on during comparison.
type Profile struct {
	Name string   `json:"name"`
	Keys []string `json:"keys"`
}

// ProfileMap maps profile names to their key sets.
type ProfileMap map[string]Profile

// LoadProfiles reads a JSON profile file and returns a ProfileMap.
func LoadProfiles(path string) (ProfileMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ProfileMap{}, nil
		}
		return nil, fmt.Errorf("read profiles: %w", err)
	}
	var profiles []Profile
	if err := json.Unmarshal(data, &profiles); err != nil {
		return nil, fmt.Errorf("parse profiles: %w", err)
	}
	pm := make(ProfileMap, len(profiles))
	for _, p := range profiles {
		pm[p.Name] = p
	}
	return pm, nil
}

// SaveProfiles writes a ProfileMap to a JSON file.
func SaveProfiles(path string, pm ProfileMap) error {
	profiles := make([]Profile, 0, len(pm))
	for _, p := range pm {
		profiles = append(profiles, p)
	}
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Name < profiles[j].Name
	})
	data, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal profiles: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// FilterByProfile returns a copy of env containing only keys listed in the profile.
// If the profile name is not found, an error is returned.
func FilterByProfile(env map[string]string, pm ProfileMap, name string) (map[string]string, error) {
	p, ok := pm[name]
	if !ok {
		return nil, fmt.Errorf("profile %q not found", name)
	}
	result := make(map[string]string, len(p.Keys))
	for _, k := range p.Keys {
		if v, exists := env[k]; exists {
			result[k] = v
		}
	}
	return result, nil
}

// ListProfileNames returns sorted profile names from a ProfileMap.
func ListProfileNames(pm ProfileMap) []string {
	names := make([]string, 0, len(pm))
	for name := range pm {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
