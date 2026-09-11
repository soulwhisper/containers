package testhelpers

import (
	"context"
	"testing"
)

// TestCNPGUserUID asserts the CloudNativePG instance-user convention: the
// postgres user is uid 26. cnpg runs instance containers (initdb included)
// as uid 26; images with a different uid fail instance bootstrap with
// 'user does not exist'.
func TestCNPGUserUID(t *testing.T, ctx context.Context, image string) {
	t.Helper()
	TestCommandSucceeds(t, ctx, image, nil, "sh", "-c", "test \"$(id -u postgres)\" = \"26\"")
}
