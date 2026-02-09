package docker

import (
	"github.com/timo-reymann/ContainerHive/internal/utils"
)

func extractTar(tarPath, destDir string) error {
	return utils.ExtractTar(tarPath, destDir)
}
