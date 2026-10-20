package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/dsh-web:latest")

	// ---- Base contract -------------------------------------------------------

	t.Run("Check dsh version", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "test \"$(dsh --version)\" = \"0.2.1-rc.2\"")
	})

	t.Run("Check entrypoint + relay exist", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "test -x /app/entrypoint.sh && test -f /app/dsh-web-relay.cjs")
	})

	// ---- Baked web profile ----------------------------------------------------

	t.Run("Check bundled versions stamp", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/opt/dsh-seed/BUNDLED_VERSIONS", nil)
	})

	t.Run("Check web profile plugins baked", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "test -d /opt/dsh-seed/profiles/web/node_modules/dsh-ui-auth && test -d /opt/dsh-seed/profiles/web/node_modules/dshmarket")
	})

	t.Run("Check default profile is web", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "test \"$DSH_PROFILE\" = web")
	})

	// ---- Runtime install/update ability (the workbench contract) --------------

	t.Run("Runtime dsh plugin add works", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image,
			&testhelpers.ContainerConfig{Env: map[string]string{"DSH_HOME": "/tmp/dshh"}},
			"sh", "-c", "dsh plugin --profile web add dsh-alive && test -d /tmp/dshh/profiles/web/node_modules/dsh-alive",
		)
	})

	// ---- Entrypoint seeding ------------------------------------------------------

	t.Run("Entrypoint seeds fresh DSH_HOME on first boot", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image,
			&testhelpers.ContainerConfig{Env: map[string]string{"DSH_HOME": "/tmp/dshh2"}},
			"sh", "-c", "timeout 15 /app/entrypoint.sh web || true; test -f /tmp/dshh2/.bundled-version && test -d /tmp/dshh2/profiles/web/node_modules/dsh-ui-auth",
		)
	})

	t.Run("Entrypoint rejects unknown profile without seeding", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image,
			&testhelpers.ContainerConfig{Env: map[string]string{
				"DSH_PROFILE": "bogus",
				"DSH_HOME":    "/tmp/dshh3",
			}},
			"sh", "-c", "! /app/entrypoint.sh && test ! -e /tmp/dshh3/.bundled-version",
		)
	})
}
