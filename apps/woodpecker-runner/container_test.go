package main

import (
	"context"
	"testing"

	"github.com/soulwhisper/containers/testhelpers"
)

func Test(t *testing.T) {
	ctx := context.Background()
	image := testhelpers.GetTestImage("ghcr.io/soulwhisper/woodpecker-runner:latest")

	tools := []struct {
		name string
		args []string
	}{
		{"crane", []string{"version"}},
		{"talosctl", []string{"version", "--client"}},
		{"kubectl", []string{"version", "--client"}},
		{"dnscontrol", []string{"version"}},
		{"shellcheck", []string{"--version"}},
		{"vendir", []string{"version"}},
		{"talosctl", []string{"version", "--client"}},
		{"yq", []string{"--version"}},
		{"flate", []string{"--help"}},
		{"jq", []string{"--version"}},
		{"git", []string{"--version"}},
		{"curl", []string{"--version"}},
		{"python3", []string{"-c", "import yaml"}},
	}

	for _, tool := range tools {
		tool := tool
		t.Run("Check tool works: "+tool.name, func(t *testing.T) {
			testhelpers.TestCommandSucceeds(t, ctx, image, nil, tool.name, tool.args...)
		})
	}
}
