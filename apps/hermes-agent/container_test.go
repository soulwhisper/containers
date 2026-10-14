package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/hermes-agent:latest")

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

	// ---- Process manager & privilege-drop tooling -------------------------

	t.Run("Check supervisord on PATH", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "which", "supervisord")
	})

	t.Run("Supervisord executes as default user", func(t *testing.T) {
		// regression: venv must not bind to a root-owned (0700) python
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "supervisord", "--version")
	})

	t.Run("Check supervisorctl on PATH", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "which", "supervisorctl")
	})

	t.Run("Check gosu on PATH", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "which", "gosu")
	})

	// ---- hermes-agent (gateway + dashboard share the same binary) ---------

	t.Run("Check hermes-agent venv python exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/opt/hermes/.venv/bin/python3", nil)
	})

	t.Run("Check hermes CLI binary exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/opt/hermes/.venv/bin/hermes", nil)
	})

	t.Run("Check feishu plugin deps import from the venv", func(t *testing.T) {
		// build-time gate: pins come from upstream's own pyproject extra; if a
		// future base drifts the closure, this fails CI before publish
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"/opt/hermes/.venv/bin/python3", "-c", "import lark_oapi.ws, requests_toolbelt, qrcode")
	})

	// ---- Baked extras ------------------------------------------------------

	t.Run("Check rtk binary exists and runs", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "/bundled/bin/rtk --version || /bundled/bin/rtk version || /bundled/bin/rtk -v")
	})

	t.Run("Check caveman skill dirs are bundled", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "for s in caveman caveman-commit caveman-review caveman-help caveman-stats caveman-compress cavecrew; do test -f \"/bundled/skills/$s/SKILL.md\" || exit 1; done")
	})

	t.Run("Check extras VERSIONS record exists", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "grep -q '^rtk=' /bundled/VERSIONS && grep -q '^caveman=' /bundled/VERSIONS && grep -q '^feishu=' /bundled/VERSIONS")
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

	// ---- Workspace mount point ---------------------------------------------

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
