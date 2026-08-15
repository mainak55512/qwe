package cli

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	in "github.com/mainak55512/qwe/initializer"
	er "github.com/mainak55512/qwe/qwerror"
	utl "github.com/mainak55512/qwe/qweutils"
	tr "github.com/mainak55512/qwe/tracker"
)

func TestHandleArgsUntrack(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := in.Init(); err != nil {
		t.Fatalf("failed to initialize repository: %v", err)
	}

	filePath := "notes.txt"
	if err := os.WriteFile(filePath, []byte("notes"), 0o644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	if _, err := tr.StartTracking(filePath); err != nil {
		t.Fatalf("failed to track file: %v", err)
	}

	originalArgs := os.Args
	t.Cleanup(func() { os.Args = originalArgs })
	os.Args = []string{"qwe", "untrack", filePath}
	if err := HandleArgs(); err != nil {
		t.Fatalf("untrack command failed: %v", err)
	}

	tracker, _, err := tr.GetTracker(tr.FileTrackerType)
	if err != nil {
		t.Fatalf("failed to read tracker: %v", err)
	}
	if _, ok := tracker[utl.Hasher(filePath)]; ok {
		t.Error("untrack command left file in tracker")
	}
	if _, err := os.Stat(filepath.Clean(filePath)); err != nil {
		t.Errorf("untrack command removed working file: %v", err)
	}

	os.Args = []string{"qwe", "untrack"}
	if err := HandleArgs(); !errors.Is(err, er.CLIUntrackErr) {
		t.Errorf("expected CLIUntrackErr, got %v", err)
	}
}
