package docker

import (
	"context"
	"errors"
	"os"

	dockerClient "github.com/docker/docker/client"
)

type Client struct {
	docker *dockerClient.Client
}

func (c *Client) Close() error {
	return c.docker.Close()
}

func NewClient() (*Client, error) {
	docker, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv, dockerClient.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Client{
		docker,
	}, nil
}

func (c *Client) LoadImageFromTar(ctx context.Context, tarPath string) (string, error) {
	imgFile, err := os.Open(tarPath)
	if err != nil {
		return "", err
	}
	defer imgFile.Close()

	res, err := c.docker.ImageLoad(ctx, imgFile)
	if err != nil {
		return "", errors.Join(errors.New("failed to load image"), err)
	}
	defer res.Body.Close()

	return imageNameFromTar(tarPath)
}
