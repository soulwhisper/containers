package testhelpers

import (
	"context"
	"fmt"
	"testing"
)

// TestUserUID asserts that the given OS user exists and has the given uid.
// The semantic check name belongs in the app's container_test.go; this is a
// generic building block.
func TestUserUID(t *testing.T, ctx context.Context, image string, user string, uid int) {
	t.Helper()
	TestCommandSucceeds(t, ctx, image, nil, "sh", "-c", fmt.Sprintf("test \"$(id -u %s)\" = \"%d\"", user, uid))
}
