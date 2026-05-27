package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/hermes-plugins:latest")

	// ---- Layout ------------------------------------------------------------

	t.Run("Check /plugins directory exists", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "test", "-d", "/plugins")
	})

	t.Run("Check /plugins is not empty", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "[ -n \"$(ls -A /plugins 2>/dev/null)\" ]")
	})

	t.Run("Check WORKDIR is /plugins", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `[ "$(pwd)" = "/plugins" ]`)
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
		testhelpers.TestFileExists(t, ctx, image, "/plugins/rtk", nil)
	})

	t.Run("Check rtk binary is executable", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "test", "-x", "/plugins/rtk")
	})

	t.Run("Check rtk binary runs", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", "/plugins/rtk --version || /plugins/rtk version || /plugins/rtk -v")
	})
}
