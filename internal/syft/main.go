package syft

import (
	"context"
	"fmt"

	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/format"
	"github.com/anchore/syft/syft/sbom"

	_ "modernc.org/sqlite" // required for rpmdb and other features
)

type SBOMImageTool struct {
	encoders *format.EncoderCollection
}

func NewSBOMImageTool() (*SBOMImageTool, error) {
	defaultEncodersConfig := format.DefaultEncodersConfig()
	encoders, err := defaultEncodersConfig.Encoders()
	if err != nil {
		return nil, err
	}

	return &SBOMImageTool{
		encoders: format.NewEncoderCollection(encoders...),
	}, nil
}

func (s *SBOMImageTool) GenerateSBOM(ctx context.Context, tarPath string) (*sbom.SBOM, error) {
	src, err := syft.GetSource(ctx, tarPath, nil)
	if err != nil {
		return nil, err
	}

	return syft.CreateSBOM(ctx, src, nil)
}

func (s *SBOMImageTool) SerializeSBOM(sbom *sbom.SBOM, outputFormat string) ([]byte, error) {
	encoder := s.encoders.GetByString(outputFormat)
	if encoder == nil {
		return nil, fmt.Errorf("unsupported output format: %s", outputFormat)
	}
	return format.Encode(*sbom, encoder)
}
