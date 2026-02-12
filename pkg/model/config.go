package model

// Secrets represents a collection of named secrets
// This is a type alias to the secrets.Secrets type
type Secrets map[string]Secret

// Secret represents a named secret with its value configuration
// This is a simplified version for the model package
type Secret struct {
	SourceType string `yaml:"source,omitempty" json:"source,omitempty" jsonschema:"Source type of the secret (env, plain). If omitted, auto-detected from value."`
	Value      string `yaml:"value" json:"value" jsonschema:"Value of the secret (env var name or plain text)"`
}

// SecretValue represents how a secret value should be resolved
type SecretValue struct {
	SourceType string `yaml:"source,omitempty" json:"source,omitempty" jsonschema:"Source type of the secret (env, plain). If omitted, auto-detected from value."`
	Value      string `yaml:"value" json:"value" jsonschema:"Value of the secret (env var name or plain text)"`
}

type VariantConfig struct {
	Name      string    `yaml:"name" json:"name" jsonschema:"Name of the variant"`
	TagSuffix string    `yaml:"tag_suffix" json:"tag_suffix" jsonschema:"Suffix to append to the tag name for this variant"`
	Versions  Versions  `yaml:"versions" json:"versions,omitempty" jsonschema:"Versions to use for this variant"`
	BuildArgs BuildArgs `yaml:"build_args" json:"build_args,omitempty" jsonschema:"Build args to add for this variant"`
}

type ImageDefinitionConfig struct {
	Tags      []*Tag          `yaml:"tags" json:"tags" jsonschema:"Tags to create for this image"`
	Variants  []VariantConfig `yaml:"variants" json:"variants,omitempty" jsonschema:"Variants to create for this image"`
	Versions  Versions        `yaml:"versions" json:"versions,omitempty" jsonschema:"Versions to use for this image"`
	BuildArgs BuildArgs       `yaml:"build_args" json:"build_args,omitempty" jsonschema:"Build args to add for this image"`
	Secrets   Secrets         `yaml:"secrets" json:"secrets,omitempty" jsonschema:"Secrets to resolve for this image"`
	DependsOn []string        `yaml:"depends_on" json:"depends_on,omitempty" jsonschema:"Names of other images in this project that must be built before this image"`
}

type HiveProjectConfig struct {
}
