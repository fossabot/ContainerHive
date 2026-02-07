package cache

type RegistryCache struct {
	CacheRef string
}

func (r RegistryCache) Name() string {
	return "registry"
}

func (r RegistryCache) ToAttributes() map[string]string {
	return map[string]string{
		"mode":           "max",
		"ref":            r.CacheRef,
		"image-manifest": "true",
		"oci-mediatypes": "true",
	}
}
