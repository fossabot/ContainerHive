package registry

import (
	"context"
	"errors"
)

// ZotRegistry is an embedded OCI registry for local development builds.
// TODO: Implement using zotregistry.dev/zot/v2 Go library.
type ZotRegistry struct{}

// NewZotRegistry creates a new ZotRegistry instance.
func NewZotRegistry() *ZotRegistry {
	return &ZotRegistry{}
}

func (z *ZotRegistry) Start(_ context.Context) error {
	return errors.New("zot registry not yet implemented")
}

func (z *ZotRegistry) Stop(_ context.Context) error {
	return nil
}

func (z *ZotRegistry) Address() string {
	return ""
}

func (z *ZotRegistry) IsLocal() bool {
	return true
}

func (z *ZotRegistry) Push(_ context.Context, _, _, _ string) error {
	return errors.New("zot registry not yet implemented")
}
