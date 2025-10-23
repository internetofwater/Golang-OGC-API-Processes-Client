// Copyright 2025 Lincoln Institute of Land Policy
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PygeoapiContainer struct {
	testcontainer testcontainers.Container
	connectionUrl string
}

// Spin up a local pygeoapi container that contains test processes
func NewPygeoapiContainer() (PygeoapiContainer, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "geopython/pygeoapi:latest",
		ExposedPorts: []string{"80/tcp"},
		WaitingFor:   wait.ForAll(wait.ForLog("Done"), wait.ForListeningPort(nat.Port("80/tcp"))),
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	genericContainer, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		return PygeoapiContainer{}, fmt.Errorf("generic container: %w", err)
	}

	mappedPort, err := genericContainer.MappedPort(ctx, "80/tcp")
	if err != nil {
		return PygeoapiContainer{}, fmt.Errorf("get api port: %w", err)
	}

	url := fmt.Sprintf("http://0.0.0.0:%d", mappedPort.Int())

	return PygeoapiContainer{
		testcontainer: genericContainer,
		connectionUrl: url,
	}, nil
}
