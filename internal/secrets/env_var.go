package secrets

import (
	"fmt"
	"os"
	"regexp"
)

const envVarResolver = "env"

var envVarRegex = regexp.MustCompile(`^\$\{?([A-Za-z_][A-Za-z0-9_]*)\}?$`)

// EnvVarResolver resolves environment variable secrets (e.g., ${VAR} or $VAR)
type EnvVarResolver struct{}

func (r *EnvVarResolver) Resolve(value string) (resolvedValue string, err error) {
	matches := envVarRegex.FindStringSubmatch(value)
	if len(matches) < 2 {
		return "", nil
	}
	name := matches[1] // Use the captured group, not the full match
	val, ok := os.LookupEnv(name)
	if !ok {
		return "", fmt.Errorf("environment variable %q not found", name)
	}
	return val, nil
}
