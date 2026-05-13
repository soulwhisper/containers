package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/hermes-suite:latest")

	// ---- Entrypoint & supervisord config ----------------------------------

	t.Run("Check entrypoint exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/app/entrypoint.sh", nil)
	})

	t.Run("Check entrypoint is executable", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "test", "-x", "/app/entrypoint.sh")
	})

	t.Run("Check entrypoint shell syntax", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "bash", "-n", "/app/entrypoint.sh")
	})

	t.Run("Check supervisord config exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/app/supervisord.conf", nil)
	})

	t.Run("Check program snippet: gateway", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/app/programs/gateway.conf", nil)
	})

	t.Run("Check program snippet: dashboard", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/app/programs/dashboard.conf", nil)
	})

	t.Run("Check program snippet: webui", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/app/programs/webui.conf", nil)
	})

	// ---- Process manager & privilege-drop tooling -------------------------

	t.Run("Check supervisord on PATH", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "which", "supervisord")
	})

	t.Run("Check supervisorctl on PATH", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "which", "supervisorctl")
	})

	t.Run("Check gosu on PATH", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "which", "gosu")
	})

	t.Run("Check bash exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/bin/bash", nil)
	})

	// ---- hermes-agent (gateway + dashboard share the same binary) ---------

	t.Run("Check hermes-agent venv python exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/opt/hermes/.venv/bin/python3", nil)
	})

	t.Run("Check hermes CLI binary exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/opt/hermes/.venv/bin/hermes", nil)
	})

	t.Run("Check tinker-atropos installed in agent venv", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"/opt/hermes/.venv/bin/python3", "-c",
			"import importlib.metadata; importlib.metadata.distribution('tinker-atropos')")
	})

	// ---- hermes-webui ------------------------------------------------------

	t.Run("Check hermes-webui server.py exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/opt/hermes-webui/server.py", nil)
	})

	t.Run("Check hermes-webui version file baked in", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/opt/hermes-webui/api/_version.py", nil)
	})

	// ---- Non-root default --------------------------------------------------

	t.Run("Default user is not root", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `[ "$(id -u)" != "0" ]`)
	})

	t.Run("Default user is hermes", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `[ "$(id -un)" = "hermes" ]`)
	})

	t.Run("Default user can write to HERMES_HOME", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "touch /opt/data/.write-test && rm /opt/data/.write-test")
	})

	// ---- Workspace mount point --------------------------------------------

	t.Run("Check workspace dir exists", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "test", "-d", "/workspace")
	})

	// ---- Entrypoint service selection --------------------------------------

	t.Run("Entrypoint rejects unknown service", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image,
			&testhelpers.ContainerConfig{Env: map[string]string{"HERMES_SERVICES": "bogus"}},
			"sh", "-c", "! /app/entrypoint.sh true",
		)
	})
}
