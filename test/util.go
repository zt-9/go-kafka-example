package test

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// container struct
type ctr struct {
	testcontainers.Container
	URI string
}

func setUpRedis(ctx context.Context) (*ctr, error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379:6379"},
		WaitingFor:   wait.ForLog("Redis Ready to accept connections").WithStartupTimeout(5*time.Minute), // default timeout is 60s
		
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "6379")

	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

	return &ctr{Container: container, URI: uri}, nil
}

func setUpKafka(ctx context.Context) (*ctr, error) {
	req := testcontainers.ContainerRequest{
		Image:        "docker.io/bitnami/kafka:3.6",
		ExposedPorts: []string{"9092/tcp"},
		WaitingFor:   wait.ForLog("Kafka Ready to accept connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "9092")

	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

	return &ctr{Container: container, URI: uri}, nil
}
