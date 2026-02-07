package cache

import "testing"

func TestRegistryCacheAttributes(t *testing.T) {
	cache := &RegistryCache{
		CacheRef: "registry.example.com/my-cache:latest",
	}

	attrs := cache.ToAttributes()

	expectedAttrs := map[string]string{
		"mode":           "max",
		"ref":            "registry.example.com/my-cache:latest",
		"image-manifest": "true",
		"oci-mediatypes": "true",
	}

	for key, want := range expectedAttrs {
		if got := attrs[key]; got != want {
			t.Errorf("attribute %q = %q, want %q", key, got, want)
		}
	}

	if cache.Name() != "registry" {
		t.Errorf("Name() = %q, want %q", cache.Name(), "registry")
	}
}
