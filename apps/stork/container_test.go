package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/stork:latest")

	t.Run("Check entrypoint exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/app/entrypoint.sh", nil)
	})

	t.Run("Check entrypoint is executable", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "test", "-x", "/app/entrypoint.sh")
	})

	t.Run("Check entrypoint shell syntax", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "sh", "-n", "/app/entrypoint.sh")
	})

	t.Run("Check supervisor log dir exists", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "test", "-d", "/var/log/supervisor")
	})

	t.Run("Check bash exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/bin/bash", nil)
	})

	t.Run("Check supervisord on PATH", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "which", "supervisord")
	})

	t.Run("Check stork-server on PATH", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "which", "stork-server")
	})

	t.Run("Check stork-agent on PATH", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "which", "stork-agent")
	})

	// Entrypoint must reject unknown STORK_MODE with non-zero exit.
	// We invert the exit code via `! cmd` so a correct rejection produces exit 0.
	t.Run("Entrypoint rejects invalid STORK_MODE", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image,
			&testhelpers.ContainerConfig{Env: map[string]string{"STORK_MODE": "invalid"}},
			"sh", "-c", "! /app/entrypoint.sh",
		)
	})

	// Entrypoint must reject a missing STORK_SUPERVISOR_CONF path.
	t.Run("Entrypoint rejects missing supervisor conf", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image,
			&testhelpers.ContainerConfig{Env: map[string]string{"STORK_SUPERVISOR_CONF": "/nonexistent.conf"}},
			"sh", "-c", "! /app/entrypoint.sh",
		)
	})
}
