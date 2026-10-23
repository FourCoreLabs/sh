// FourCoreLabs fork patch; see FORK.md.

package interp_test

import (
	"bytes"
	"context"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/FourCoreLabs/sh/interp"
	"github.com/FourCoreLabs/sh/syntax"
)

// boomHandler panics when the command "boom" runs, standing in for a
// builtin or handler bug in third-party code.
func boomHandler(next interp.ExecHandlerFunc) interp.ExecHandlerFunc {
	return func(ctx context.Context, args []string) error {
		if args[0] == "boom" {
			panic("boom handler bug")
		}
		return next(ctx, args)
	}
}

// TestPanicInGoroutineStageIsContained: a panic in a pipeline stage,
// background job or process substitution runs on a goroutine the caller
// cannot recover. Without recoverStage it kills the host process. With it,
// the stage fails with status 2 and the script carries on.
func TestPanicInGoroutineStageIsContained(t *testing.T) {
	tests := []struct {
		name, script, wantOut string
	}{
		{"pipeline left stage", "boom | true; echo after", "after\n"},
		{"background job", "boom & wait $!; echo status=$?", "status=2\n"},
		{"process substitution", "true <(boom); echo after", "after\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "process substitution" && runtime.GOOS == "windows" {
				t.Skip("process substitution is not supported on Windows upstream")
			}
			file, err := syntax.NewParser().Parse(strings.NewReader(tt.script), "")
			if err != nil {
				t.Fatal(err)
			}
			var stdout, stderr bytes.Buffer
			r, err := interp.New(
				interp.StdIO(nil, &stdout, &stderr),
				interp.ExecHandlers(boomHandler),
			)
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_ = r.Run(ctx, file)
			if stdout.String() != tt.wantOut {
				t.Errorf("stdout = %q, want %q (stderr %q)", stdout.String(), tt.wantOut, stderr.String())
			}
			if !strings.Contains(stderr.String(), "internal error: boom handler bug") {
				t.Errorf("stderr = %q, want the internal error diagnostic", stderr.String())
			}
		})
	}
}
