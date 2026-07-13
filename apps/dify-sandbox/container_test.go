package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/dify-sandbox:latest")

	// ---- sandbox user exists ----------------------------------------------

	t.Run("sandbox user exists", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"id", "sandbox")
	})

	t.Run("sandbox user has bash login shell", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `[ "$(getent passwd sandbox | cut -d: -f7)" = "/bin/bash" ]`)
	})

	t.Run("sandbox home is writable", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"su", "-", "sandbox", "-c", `touch "$HOME/.probe"`)
	})

	// ---- mise on system PATH ----------------------------------------------

	t.Run("which mise", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "which", "mise")
	})

	t.Run("mise runs", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "mise", "--version")
	})

	// ---- node (built-in via mise) -----------------------------------------

	t.Run("which node", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "which", "node")
	})

	t.Run("node runs", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "node", "--version")
	})

	// ---- officecli (npm package installed via mise) -----------------------

	t.Run("which officecli", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "which", "officecli")
	})

	// ---- tools available to sandbox user ----------------------------------

	t.Run("sandbox user can run node", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"su", "-", "sandbox", "-c", "node --version")
	})

	t.Run("sandbox user can run officecli", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"su", "-", "sandbox", "-c", "which officecli")
	})
}
