package file_resolver

import "fmt"

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
		for _, tmplExt := range SupportedTemplateExtensions {
			for _, ext := range extensions {
				possibleNames[idx] = fmt.Sprintf("%s.%s.%s", baseName, ext, tmplExt)
				idx++
			}
		}
	}

	return possibleNames
}
