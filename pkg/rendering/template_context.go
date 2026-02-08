package rendering

import (
	"github.com/timo-reymann/ContainerHive/internal/buildconfig_resolver"
	"github.com/timo-reymann/ContainerHive/internal/file_resolver/templating"
	"github.com/timo-reymann/ContainerHive/pkg/model"
)

func newTemplateContext(image *model.Image, values *buildconfig_resolver.ResolvedBuildValues) *templating.TemplateContext {
	return &templating.TemplateContext{
		ImageName: image.Name,
		Versions:  values.Versions,
		BuildArgs: values.BuildArgs,
	}
}
