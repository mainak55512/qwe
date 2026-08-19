package binaryhandler

import (
	"os"
	"path/filepath"
	"testing"

	cp "github.com/mainak55512/qwe/compressor"
)

func TestRevertBinFileReplacesTargetAtomically(t *testing.T) {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(originalDir) })

	objectDir := filepath.Join(".qwe", "_object")
	if err := os.MkdirAll(objectDir, 0o755); err != nil {
		t.Fatal(err)
	}
	object := filepath.Join(objectDir, "revision")
	if err := os.WriteFile(object, []byte{0, 1, 2, 3}, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := cp.CompressFile(object); err != nil {
		t.Fatal(err)
	}

	target := filepath.Join(t.TempDir(), "binary")
	if err := os.WriteFile(target, []byte("working data"), 0o640); err != nil {
		t.Fatal(err)
	}
	if err := RevertBinFile(target, "revision"); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != string([]byte{0, 1, 2, 3}) {
		t.Fatalf("unexpected restored content: %v", content)
	}
	info, err := os.Stat(target)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o640 {
		t.Fatalf("target permissions changed: %o", info.Mode().Perm())
	}
}
