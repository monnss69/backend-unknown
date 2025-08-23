package components

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
)

const maxCodeSize = 200 * 1024

var importRE = regexp.MustCompile(`(?m)^\s*import\s+`)

// validateCode ensures the component code meets basic constraints.
func validateCode(code string) error {
	if len(code) > maxCodeSize {
		return errors.New("code exceeds 200KB limit")
	}
	if importRE.MatchString(code) {
		return errors.New("imports are not allowed")
	}
	return nil
}

// extractPropsSchema attempts to derive a simple props schema from the code.
// It looks for an `interface Props { ... }` declaration.
func extractPropsSchema(code string) map[string]string {
	re := regexp.MustCompile(`interface\s+Props\s*{([^}]*)}`)
	match := re.FindStringSubmatch(code)
	if len(match) < 2 {
		return map[string]string{}
	}
	body := match[1]
	schema := make(map[string]string)
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}
		name := strings.TrimSpace(strings.TrimSuffix(parts[0], "?"))
		typ := strings.TrimSpace(strings.TrimSuffix(parts[1], ";"))
		schema[name] = typ
	}
	return schema
}

// randomID returns a random hex string suitable for database keys.
func randomID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
