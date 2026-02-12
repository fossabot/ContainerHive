package secrets

const plainTextResolver = "plain"

// PlainTextResolver resolves plain text secrets (returns the value as-is)
// Only handles values that don't match environment variable patterns
type PlainTextResolver struct{}

func (r *PlainTextResolver) Resolve(value string) (resolvedValue string, err error) {
	return value, nil
}
