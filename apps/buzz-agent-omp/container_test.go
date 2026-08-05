package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/buzz-agent-omp:latest")

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

	t.Run("HOME is the agent home and writable", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `[ "$HOME" = "/home/agent" ] && touch "$HOME/.probe"`)
	})

	// ---- Tooling on PATH (sprig multicall + mise shims) --------------------

	for _, bin := range []string{
		"buzz", "buzz-acp", "buzz-agent", "buzz-dev-mcp",
		"git-credential-nostr", "git-sign-nostr",
		"gh", "git", "jq", "just", "mise", "omp", "prek", "rg", "rtk", "tmux", "tree", "uv", "yq",
	} {
		bin := bin
		t.Run("which "+bin, func(t *testing.T) {
			testhelpers.TestCommandSucceeds(t, ctx, image, nil, "which", bin)
		})
	}

	t.Run("git runs", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "git", "--version")
	})

	t.Run("omp runs", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "omp", "--version")
	})

	t.Run("buzz-acp runs", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil, "buzz-acp", "--help")
	})

	// ---- git configuration (workspace + nostr signing contract) ------------

	t.Run("git safe.directory is wildcarded", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `git config --system --get-all safe.directory | grep -qx '*'`)
	})

	t.Run("git signs commits via git-sign-nostr (x509)", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `[ "$(git config --system gpg.format)" = "x509" ] && [ "$(git config --system gpg.x509.program)" = "/usr/local/bin/git-sign-nostr" ] && [ "$(git config --system commit.gpgSign)" = "true" ]`)
	})

	// ---- Harness defaults: omp is the baked ACP agent ----------------------

	t.Run("omp is the default ACP agent command", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `[ "$BUZZ_ACP_AGENT_COMMAND" = "omp" ] && [ "$BUZZ_ACP_AGENT_ARGS" = "acp" ]`)
	})

	t.Run("Agent command override passes through from container env", func(t *testing.T) {
		cfg := &testhelpers.ContainerConfig{
			Env: map[string]string{"BUZZ_ACP_AGENT_COMMAND": "buzz-agent"},
		}
		testhelpers.TestCommandSucceeds(t, ctx, image, cfg,
			"sh", "-c", `[ "$BUZZ_ACP_AGENT_COMMAND" = "buzz-agent" ]`)
	})

	// ---- Identity / relay routing (env-driven, fail closed) ----------------

	t.Run("No agent identity baked by default", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `[ -z "${BUZZ_PRIVATE_KEY:-}" ]`)
	})

	t.Run("Entrypoint refuses to launch without an identity key", func(t *testing.T) {
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `! /app/entrypoint.sh 2>/dev/null`)
	})

	t.Run("Entrypoint wires relay-scoped git credential helper", func(t *testing.T) {
		// Invalid nsec fails fast in buzz-acp after the git config is written;
		// the assertion checks the helper wiring, not the launch.
		testhelpers.TestCommandSucceeds(t, ctx, image, nil,
			"sh", "-c", `BUZZ_PRIVATE_KEY=invalid BUZZ_RELAY_URL=wss://relay.example /app/entrypoint.sh 2>/dev/null; [ "$(git config --global --get 'credential.https://relay.example/git.helper')" = "/usr/local/bin/git-credential-nostr" ] && [ "$(git config --global --get 'credential.https://relay.example/git.useHttpPath')" = "true" ]`)
	})
}
