package main

import (
	"context"
	"testing"

	"github.com/soulwhiser/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/forgejo-runner:latest")

	t.Run("Check cosign exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/usr/local/bin/cosign", nil)
	})

	t.Run("Check flux exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/usr/local/bin/flux", nil)
	})

	t.Run("Check flux-local exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/root/.local/bin/flux-local", nil)
	})

	t.Run("Check helm exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/usr/local/bin/helm", nil)
	})

	t.Run("Check kustomize exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/usr/local/bin/kustomize", nil)
	})

	t.Run("Check rsync exists", func(t *testing.T) {
		testhelpers.TestFileExists(t, ctx, image, "/usr/bin/rsync", nil)
	})
}
