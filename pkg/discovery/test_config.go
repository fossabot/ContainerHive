package discovery

import (
	"errors"
	"os"
	"path/filepath"
)

var testConfigFileNames = []string{"test.yaml", "test.yml", "test.yml.gotpl", "test.yaml.gotpl"}

func getTestConfigFilePath(root string) (string, error) {
	for _, name := range testConfigFileNames {
		_, err := os.Stat(filepath.Join(root, name))
		if err != nil && !os.IsNotExist(err) {
			return "", errors.Join(errors.New("failed to stat test config file path "+name), err)
		}
		if err == nil {
			return filepath.Join(root, name), nil
		}
	}
	return "", nil
}
