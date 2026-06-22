package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/devbox:latest")

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

	t.Run("HOME is the dev home and writable", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `[ "$HOME" = "/home/dev" ] && touch "$HOME/.probe"`)
	})

	// ---- Tooling on PATH (mise shims) -------------------------------------

	for _, bin := range []string{"mise", "git", "node", "claude", "gh", "just", "prek", "jq", "yq"} {
		bin := bin
		t.Run("which "+bin, func(t *testing.T) {
			testhelpers.TestCommandSucceeds(t, ctx, image, nil, "which", bin)
		})
	}

	t.Run("git runs", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "git", "--version")
	})

	t.Run("claude runs", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "claude", "--version")
	})

	t.Run("git safe.directory is wildcarded", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `git config --system --get-all safe.directory | grep -qx '*'`)
	})

	// ---- Provider-neutral: nothing baked, env drives everything -----------

	t.Run("No provider endpoint baked by default", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `[ -z "${ANTHROPIC_BASE_URL:-}" ]`)
	})

	t.Run("Routing endpoint passes through from container env", func(t *testing.T) {
		cfg := &testhelpers.ContainerConfig{
			Env: map[string]string{"ANTHROPIC_BASE_URL": "https://api.deepseek.com/anthropic"},
		}
		testhelpers.TestCommandSucceeds(t, ctx, image, cfg,
			"sh", "-c", `[ "$ANTHROPIC_BASE_URL" = "https://api.deepseek.com/anthropic" ]`)
	})

	t.Run("Auth token passthrough is visible to the process", func(t *testing.T) {
		cfg := &testhelpers.ContainerConfig{
			Env: map[string]string{"ANTHROPIC_AUTH_TOKEN": "sk-probe"},
		}
		testhelpers.TestCommandSucceeds(t, ctx, image, cfg,
			"sh", "-c", `[ "$ANTHROPIC_AUTH_TOKEN" = "sk-probe" ]`)
	})
}
