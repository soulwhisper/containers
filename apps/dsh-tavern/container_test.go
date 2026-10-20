package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/dsh-tavern:latest")

	// ---- Base contract -------------------------------------------------------

	t.Run("Check dsh version", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "test \"$(dsh --version)\" = \"0.2.1-rc.2\"")
	})

	t.Run("Check entrypoint + relay exist", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "test -x /app/entrypoint.sh && test -f /app/dsh-web-relay.cjs")
	})

	t.Run("Check default profile is tavern", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "test \"$DSH_PROFILE\" = tavern")
	})

	// ---- Recipe artifacts -----------------------------------------------------

	t.Run("Check tavern source with workspace deps", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"test", "-d", "/opt/dsh-tavern/apps/dsh-tavern/node_modules/jsdom")
	})

	t.Run("Check bundled versions stamp", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/opt/dsh-seed/BUNDLED_VERSIONS", nil)
	})

	t.Run("Check profile plugin dirs baked", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "cd /opt/dsh-seed/profiles/tavern/node_modules && test -d dsh-tavern-wrongbook && test -d dsh-tavern-comfy && test -d dsh-tavern-knock-knock && test -d narrative-anchor")
	})

	t.Run("Check font-limit host patch applied", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "grep -rq 'max: 32' /usr/local/lib/node_modules/@deepseek-ai/dsh/node_modules/@deepseek-ai/dsh-client-ui-theme/lib/")
	})

	// ---- Freeze contract --------------------------------------------------------

	t.Run("Check pnpm binaries removed (frozen runtime)", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "! command -v pnpm && ! command -v pnpx")
	})

	// ---- Boot smoke ---------------------------------------------------------------

	t.Run("Boot: tavern profile imports cleanly and activates plugins", func(t *testing.T) {
		// full boot (~25s): entrypoint seeds, relay binds, dsh imports the
		// tavern plugin (the #300 failure point) and wrongbook activates.
		// "skipping profile bundle" must NOT appear: pnpm's relative link:
		// symlinks dangle wherever the profile is copied — the Dockerfile
		// rewrites them as absolute symlinks (step 7)
		testhelpers.TestCommandSucceeds(t, ctx, image,
			&testhelpers.ContainerConfig{Env: map[string]string{"DSH_HOME": "/tmp/dshh"}},
			"sh", "-c", "timeout 45 /app/entrypoint.sh tavern > /tmp/boot.log 2>&1 || true; test \"$(grep -c 'failed to import' /tmp/boot.log)\" = 0 && ! grep -q 'skipping profile bundle' /tmp/boot.log && grep -q 'relay:' /tmp/boot.log && grep -q 'dsh-wrongbook' /tmp/boot.log && test -f /tmp/dshh/.bundled-version",
		)
	})

	t.Run("Entrypoint rejects unknown profile without seeding", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image,
			&testhelpers.ContainerConfig{Env: map[string]string{
				"DSH_PROFILE": "bogus",
				"DSH_HOME":    "/tmp/dshh2",
			}},
			"sh", "-c", "! /app/entrypoint.sh && test ! -e /tmp/dshh2/.bundled-version",
		)
	})
}
