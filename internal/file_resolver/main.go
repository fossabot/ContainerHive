package file_resolver

import (
	"fmt"
	"os"
	"path/filepath"
)

const TemplateExtensionGoTemplate = "gotpl"

var SupportedTemplateExtensions = []string{
	TemplateExtensionGoTemplate,
}

func GetFileCandidates(baseName string, extensions ...string) []string {
	extLen := len(extensions)
	var possibleNames []string

	if extLen == 0 {
		possibleNames = make([]string, len(SupportedTemplateExtensions))
		for idx, tmplExt := range SupportedTemplateExtensions {
			possibleNames[idx] = fmt.Sprintf("%s.%s", baseName, tmplExt)
		}
	} else {
		possibleNames = make([]string, extLen*len(SupportedTemplateExtensions))
		idx := 0
		for _, ext := range extensions {
			for _, tmplExt := range SupportedTemplateExtensions {
				possibleNames[idx] = fmt.Sprintf("%s.%s.%s", baseName, ext, tmplExt)
				idx++
			}
		}
	}

	return possibleNames
}

func ResolveFirstExistingFile(root string, candidates ...string) (string, error) {
	for _, candidate := range candidates {
		candidatePath := filepath.Join(root, candidate)
		if stat, err := os.Stat(candidatePath); err == nil && !stat.IsDir() {
			return candidatePath, nil
		}
	}
	return "", nil
}
