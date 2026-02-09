package registry

import (
	"context"
	"errors"
	"os"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/layout"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/timo-reymann/ContainerHive/internal/utils"
)

// RemoteRegistry is a passthrough registry for CI environments.
// Start and Stop are no-ops; Push pushes to the configured remote registry.
type RemoteRegistry struct {
	address string
}

// NewRemoteRegistry creates a remote registry with the given address.
func NewRemoteRegistry(address string) *RemoteRegistry {
	return &RemoteRegistry{address: address}
}

func (r *RemoteRegistry) Start(_ context.Context) error {
	return nil
}

func (r *RemoteRegistry) Stop(_ context.Context) error {
	return nil
}

func (r *RemoteRegistry) Address() string {
	return r.address
}

func (r *RemoteRegistry) IsLocal() bool {
	return false
}

func (r *RemoteRegistry) Push(_ context.Context, imageName, tag, ociTarPath string) error {
	tmpDir, err := os.MkdirTemp("", "oci-push-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	if err := utils.ExtractTar(ociTarPath, tmpDir); err != nil {
		return errors.Join(errors.New("failed to extract OCI tar for push"), err)
	}

	layoutPath, err := layout.FromPath(tmpDir)
	if err != nil {
		return errors.Join(errors.New("failed to read OCI layout"), err)
	}

	idx, err := layoutPath.ImageIndex()
	if err != nil {
		return err
	}

	idxManifest, err := idx.IndexManifest()
	if err != nil {
		return err
	}

	if len(idxManifest.Manifests) == 0 {
		return errors.New("no manifests in OCI layout")
	}

	img, err := layoutPath.Image(idxManifest.Manifests[0].Digest)
	if err != nil {
		return errors.Join(errors.New("failed to read image from layout"), err)
	}

	ref, err := name.NewTag(r.address + "/" + imageName + ":" + tag)
	if err != nil {
		return errors.Join(errors.New("invalid image reference"), err)
	}

	if err := remote.Write(ref, img); err != nil {
		return errors.Join(errors.New("failed to push image to remote registry"), err)
	}

	return nil
}
