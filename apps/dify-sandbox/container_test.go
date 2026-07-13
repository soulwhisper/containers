package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/dify-sandbox:latest")

	// ---- Non-root default (uid=gid=2000) ----------------------------------

	t.Run("Default user is uid 2000", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `[ "$(id -u)" = "2000" ]`)
	})

	t.Run("Default group is gid 2000", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `[ "$(id -g)" = "2000" ]`)
	})

	t.Run("Default user is not root", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `[ "$(id -u)" != "0" ]`)
	})

	t.Run("Login shell is bash", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `[ "$(getent passwd "$(id -u)" | cut -d: -f7)" = "/bin/bash" ]`)
	})

	t.Run("HOME is the sandbox home and writable", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `[ "$HOME" = "/home/sandbox" ] && touch "$HOME/.probe"`)
	})

	// ---- mise on PATH ----------------------------------------------------

	t.Run("which mise", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "which", "mise")
	})

	t.Run("mise runs", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "mise", "--version")
	})

	// ---- node (built-in via mise) ----------------------------------------

	t.Run("which node", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "which", "node")
	})

	t.Run("node runs", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "node", "--version")
	})

	// ---- officecli (npm package installed via mise) ----------------------

	t.Run("which officecli", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "which", "officecli")
	})

	t.Run("officecli runs", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "officecli", "--help")
	})
}
