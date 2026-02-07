package container_structure_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/GoogleContainerTools/container-structure-test/cmd/container-structure-test/app/cmd/test"
	"github.com/GoogleContainerTools/container-structure-test/pkg/config"
	"github.com/GoogleContainerTools/container-structure-test/pkg/drivers"
	"github.com/GoogleContainerTools/container-structure-test/pkg/types/unversioned"
	docker "github.com/docker/docker/client"
)

type TestRunner struct {
	TestDefinitionPath string
	ImageTarPath       string
	Platform           string
	ReportFile         string
}

func (t *TestRunner) getOptions(output unversioned.OutputValue) *config.StructureTestOptions {
	return &config.StructureTestOptions{
		ImagePath:           t.ImageTarPath,
		IgnoreRefAnnotation: false,
		ConfigFiles:         []string{t.TestDefinitionPath},
		Platform:            t.Platform,
		JSON:                true,
		Output:              output,
		NoColor:             false,
		Driver:              "docker",
		Quiet:               true,
	}
}

func (t *TestRunner) isTar() bool {
	return filepath.Ext(t.ImageTarPath) == ".tar"
}

func (t *TestRunner) loadImageFromTar(ctx context.Context) (string, error) {
	c, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithHostFromEnv())
	if err != nil {
		return "", err
	}
	defer c.Close()

	imgFile, err := os.Open(t.ImageTarPath)
	if err != nil {
		return "", err
	}
	defer imgFile.Close()

	res, err := c.ImageLoad(ctx, imgFile)
	if err != nil {
		return "", errors.Join(errors.New("failed to load image"), err)
	}
	defer res.Body.Close()

	return imageNameFromTar(t.ImageTarPath)
}

func (t *TestRunner) resolveImageName(ctx context.Context) (string, error) {
	if t.isTar() {
		return t.loadImageFromTar(ctx)
	}
	return t.ImageTarPath, nil
}

func (t *TestRunner) runTests(channel chan interface{}, imageName string, opts *config.StructureTestOptions) {
	args := &drivers.DriverConfig{
		Image:    imageName,
		Save:     opts.Save,
		Metadata: opts.Metadata,
		Runtime:  opts.Runtime,
		Platform: opts.Platform,
	}
	driverImpl := drivers.InitDriverImpl(opts.Driver)
	tests, err := test.Parse(t.TestDefinitionPath, args, driverImpl)
	if err != nil {
		channel <- &unversioned.TestResult{
			Errors: []string{
				fmt.Sprintf("error parsing config file: %s", err),
			},
		}
	}
	tests.RunAll(channel, t.TestDefinitionPath)
	close(channel)
}

func (t *TestRunner) Run() error {
	imageName, err := t.resolveImageName(context.Background())
	if err != nil {
		return err
	}

	opts := t.getOptions(unversioned.Junit)
	channel := make(chan interface{}, 1)
	go t.runTests(channel, imageName, opts)

	testReportFile, err := os.Create(t.ReportFile)
	if err != nil {
		return err
	}
	defer testReportFile.Close()

	return test.ProcessResults(testReportFile, unversioned.Junit, opts.JunitSuiteName, channel)
}
