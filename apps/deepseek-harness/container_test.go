package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/deepseek-harness:latest")

	// ---- Entrypoint & relay -----------------------------------------------

	t.Run("Check entrypoint exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/app/entrypoint.sh", nil)
	})

	t.Run("Check entrypoint is executable", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "test", "-x", "/app/entrypoint.sh")
	})

	t.Run("Check entrypoint shell syntax", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "bash", "-n", "/app/entrypoint.sh")
	})

	t.Run("Check web relay script exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/app/dsh-web-relay.cjs", nil)
	})

	// ---- Core tooling -------------------------------------------------------

	t.Run("Check dsh version", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "test \"$(dsh --version)\" = \"0.2.1-rc.2\"")
	})

	t.Run("Check pnpm on PATH", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "pnpm", "--version")
	})

	t.Run("Check tini exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/usr/bin/tini", nil)
	})

	t.Run("Check runs as node user", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "test \"$(id -un)\" = node")
	})

	// ---- Baked plugin seed --------------------------------------------------

	t.Run("Check bundled versions stamp", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/opt/dsh-seed/BUNDLED_VERSIONS", nil)
	})

	t.Run("Check web profile plugins baked", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "test -d /opt/dsh-seed/profiles/web/node_modules/dsh-ui-auth && test -d /opt/dsh-seed/profiles/web/node_modules/dsh-pocket-relay && test -d /opt/dsh-seed/profiles/web/node_modules/dshmarket")
	})

	t.Run("Check tavern profile plugins baked", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "test -d /opt/dsh-seed/profiles/tavern/node_modules/dsh-ui-auth && test -d /opt/dsh-seed/profiles/tavern/node_modules/dsh-pocket")
	})

	t.Run("Check tavern runtime present in seed", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "ls /opt/dsh-seed/profiles/tavern/node_modules | grep -qi tavern")
	})

	// ---- Entrypoint profile handling ---------------------------------------

	t.Run("Entrypoint rejects unknown profile without seeding", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image,
			&testhelpers.ContainerConfig{Env: map[string]string{
				"DSH_PROFILE": "bogus",
				"DSH_HOME":    "/tmp/dshh",
			}},
			"sh", "-c", "! /app/entrypoint.sh && test ! -e /tmp/dshh/.bundled-version",
		)
	})

	t.Run("Entrypoint seeds fresh DSH_HOME on first boot", func(t *testing.T) {
		// dsh web may still be starting when timeout kills it (124); the seed
		// side effects are what we assert, not web readiness
		testhelpers.TestCommandSucceeds(t, ctx, image,
			&testhelpers.ContainerConfig{Env: map[string]string{"DSH_HOME": "/tmp/dshh"}},
			"sh", "-c", "timeout 15 /app/entrypoint.sh web || true; test -f /tmp/dshh/.bundled-version && test -d /tmp/dshh/profiles/web/node_modules/dsh-ui-auth && test -d /tmp/dshh/profiles/tavern/node_modules/dsh-pocket",
		)
	})
}
