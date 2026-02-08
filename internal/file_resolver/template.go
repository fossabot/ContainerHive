package file_resolver

import (
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/fileutils"
)

func CopyAndRenderFile(sourceFile, targetPath string) error {
	ext, _ := strings.CutPrefix(filepath.Ext(sourceFile), ".")

	switch ext {
	case TemplateExtensionGoTemplate:
		// TODO Template
		break
	}

	_, err := fileutils.CopyFile(sourceFile, targetPath)
	return err
}
