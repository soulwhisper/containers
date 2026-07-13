package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/dify-sandbox:latest")

	// ---- upstream image is intact -----------------------------------------

	t.Run("upstream /main binary exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/main", nil)
	})

	t.Run("upstream /entrypoint.sh exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/entrypoint.sh", nil)
	})

	// ---- mise on system PATH ----------------------------------------------

	t.Run("command -v mise", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "command -v mise")
	})

	t.Run("mise runs", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "mise", "--version")
	})

	// ---- node (built-in via mise) -----------------------------------------

	t.Run("command -v node", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "command -v node")
	})

	t.Run("node runs", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "node", "--version")
	})

	// ---- officecli (npm package installed via mise) -----------------------

	t.Run("command -v officecli", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "command -v officecli")
	})

	// ---- sandbox user exists ----------------------------------------------

	t.Run("sandbox user exists", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "id", "sandbox")
	})
}
