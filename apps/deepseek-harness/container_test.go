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

	t.Run("Check pnpm version", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "test \"$(pnpm --version)\" = \"11.25.0\"")
	})

	t.Run("Check tini exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/usr/bin/tini", nil)
	})

	t.Run("Check runs as node user", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "test \"$(id -un)\" = node")
	})

	// ---- Basement contract --------------------------------------------------

	t.Run("Check no plugin seed shipped in base", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "test ! -e /opt/dsh-seed")
	})

	t.Run("Entrypoint rejects unknown profile without seeding", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image,
			&testhelpers.ContainerConfig{Env: map[string]string{
				"DSH_PROFILE": "bogus",
				"DSH_HOME":    "/tmp/dshh",
			}},
			"sh", "-c", "! /app/entrypoint.sh && test ! -e /tmp/dshh/.bundled-version",
		)
	})

	t.Run("Entrypoint runs stock dsh web without a seed", func(t *testing.T) {
		// the relay log line proves the entrypoint reached the run stage;
		// dsh itself may still be starting when timeout kills it
		testhelpers.TestCommandSucceeds(t, ctx, image,
			&testhelpers.ContainerConfig{Env: map[string]string{"DSH_HOME": "/tmp/dshh"}},
			"sh", "-c", "timeout 10 /app/entrypoint.sh web 2>&1 | grep -q \"relay:\"",
		)
	})

	t.Run("Entrypoint passes DSH_TRUSTED_HOSTS/DSH_PUBLIC_URL through", func(t *testing.T) {
		// with --public-url the startup line prints the advertised root;
		// that only happens if the app actually received both flags
		testhelpers.TestCommandSucceeds(t, ctx, image,
			&testhelpers.ContainerConfig{Env: map[string]string{
				"DSH_HOME":          "/tmp/dshh2",
				"DSH_TRUSTED_HOSTS": "tavern.example,192.168.1.10:9400",
				"DSH_PUBLIC_URL":    "https://tavern.example",
			}},
			"sh", "-c", "timeout 12 /app/entrypoint.sh web 2>&1 | grep -q \"https://tavern.example/\"",
		)
	})
}
