package testutil

import "runtime/debug"

const buildkitModule = "github.com/moby/buildkit"

// BuildKitImage returns the moby/buildkit Docker image tag matching the version in go.mod.
func BuildKitImage() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "moby/buildkit:latest"
	}
	for _, dep := range info.Deps {
		if dep.Path == buildkitModule {
			return "moby/buildkit:" + dep.Version
		}
	}
	return "moby/buildkit:latest"
}
