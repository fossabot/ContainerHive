package docker

import (
	"archive/tar"
	"encoding/json"
	"errors"
	"os"
)

func imageNameFromTar(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", errors.Join(errors.New("failed to open image tar"), err)
	}
	defer f.Close()

	tr := tar.NewReader(f)
	for {
		hdr, err := tr.Next()
		if err != nil {
			return "", errors.Join(errors.New("index.json not found in image tar"), err)
		}
		if hdr.Name == "index.json" {
			break
		}
	}

	var index struct {
		Manifests []struct {
			Annotations map[string]string `json:"annotations"`
		} `json:"manifests"`
	}
	if err := json.NewDecoder(tr).Decode(&index); err != nil {
		return "", errors.Join(errors.New("failed to decode index.json"), err)
	}
	if len(index.Manifests) == 0 {
		return "", errors.New("no manifests found in index.json")
	}
	imageName, ok := index.Manifests[0].Annotations["io.containerd.image.name"]
	if !ok {
		return "", errors.New("no image name annotation found in index.json")
	}

	return imageName, nil
}
