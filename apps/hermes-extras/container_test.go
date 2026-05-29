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
}
