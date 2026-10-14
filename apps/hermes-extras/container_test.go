package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/hermes-extras:latest")

	// ---- Plugin python deps (the only payload) ------------------------------

	t.Run("Check plugin deps are bundled", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/pylib/lark_oapi/__init__.py", nil)
		testhelpers.TestFileExists(t, ctx, image, "/pylib/requests_toolbelt/__init__.py", nil)
	})

	// ---- Non-root default ---------------------------------------------------

	t.Run("Default user is uid 10000", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `[ "$(id -u)" = "10000" ]`)
	})

	// ---- install.sh (initContainer entrypoint) ------------------------------

	t.Run("Check install.sh mirrors plugin deps", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image,
			&testhelpers.ContainerConfig{Env: map[string]string{
				"PYLIB_DIR": "/tmp/pylib",
			}},
			"sh", "-c", "/app/install.sh && test -f /tmp/pylib/lark_oapi/__init__.py && test -f /tmp/pylib/requests_toolbelt/__init__.py")
	})

	t.Run("Check install.sh is idempotent and evicts stale pins", func(t *testing.T) {
		// a foreign leftover from a previous pin must be removed by the mirror,
		// or two lark_oapi copies would fight on the same PYTHONPATH entry
		testhelpers.TestCommandSucceeds(t, ctx, image,
			&testhelpers.ContainerConfig{Env: map[string]string{
				"PYLIB_DIR": "/tmp/pylib2",
			}},
			"sh", "-c", "mkdir -p /tmp/pylib2/lark_oapi-9.9.9.dist-info && /app/install.sh && test ! -e /tmp/pylib2/lark_oapi-9.9.9.dist-info && /app/install.sh && test -f /tmp/pylib2/lark_oapi/__init__.py")
	})
}
