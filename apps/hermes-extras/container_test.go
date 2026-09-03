package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/hermes-extras:latest")

	// ---- Layout ------------------------------------------------------------

	t.Run("Check /data directory exists", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "test", "-d", "/data")
	})

	t.Run("Check /data is not empty", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "[ -n \"$(ls -A /data 2>/dev/null)\" ]")
	})

	t.Run("Check WORKDIR is /data", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `[ "$(pwd)" = "/data" ]`)
	})

	// ---- Build-time tooling still present (debugability) ------------------

	t.Run("Check mise on PATH", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "which", "mise")
	})

	// ---- Non-root default --------------------------------------------------

	t.Run("Default user is uid 10000", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `[ "$(id -u)" = "10000" ]`)
	})

	t.Run("Default user is not root", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `[ "$(id -u)" != "0" ]`)
	})

	// ---- Bundled plugins ---------------------------------------------------

	t.Run("Check rtk binary exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/data/rtk", nil)
	})

	t.Run("Check rtk binary is executable", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "test", "-x", "/data/rtk")
	})

	t.Run("Check rtk binary runs", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "/data/rtk --version || /data/rtk version || /data/rtk -v")
	})

	t.Run("Check buzz binary exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/data/buzz", nil)
	})

	t.Run("Check buzz binary is executable", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "test", "-x", "/data/buzz")
	})

	t.Run("Check buzz binary runs", func(t *testing.T) {
		// buzz-cli has no clap `version` attribute — `--version` is an unknown
		// flag (exit 1). `--help` is the supported zero-side-effect success path.
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "/data/buzz", "--help")
	})

	// ---- caveman skills ----------------------------------------------------

	t.Run("Check all caveman skill dirs are bundled", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "for s in caveman caveman-commit caveman-review caveman-help caveman-stats caveman-compress cavecrew; do test -f \"/skills/$s/SKILL.md\" || exit 1; done")
	})

	// ---- install.sh (initContainer entrypoint) -----------------------------

	t.Run("Check install.sh mirrors binaries and skills", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image,
			&testhelpers.ContainerConfig{Env: map[string]string{
				"BIN_DIR":    "/tmp/bin",
				"SKILLS_DIR": "/tmp/skills",
			}},
			"sh", "-c", "/app/install.sh && test -x /tmp/bin/rtk && test -x /tmp/bin/buzz && test -f /tmp/skills/caveman/SKILL.md")
	})

	t.Run("Check install.sh is idempotent", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image,
			&testhelpers.ContainerConfig{Env: map[string]string{
				"BIN_DIR":    "/tmp/bin",
				"SKILLS_DIR": "/tmp/skills",
			}},
			"sh", "-c", "/app/install.sh && /app/install.sh && test -f /tmp/skills/cavecrew/SKILL.md")
	})
}
