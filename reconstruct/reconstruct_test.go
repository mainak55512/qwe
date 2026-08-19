package reconstruct

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	cp "github.com/mainak55512/qwe/compressor"
	tr "github.com/mainak55512/qwe/tracker"
)

func setupRepository(t *testing.T) string {
	t.Helper()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(originalDir) })
	if err := os.MkdirAll(filepath.Join(".qwe", "_object"), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func writeCompressedObject(t *testing.T, name string, content []byte) {
	t.Helper()
	path := filepath.Join(".qwe", "_object", name)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := cp.CompressFile(path); err != nil {
		t.Fatal(err)
	}
}

func TestReconstructLeavesTargetUnchangedWhenRevisionIsCorrupt(t *testing.T) {
	setupRepository(t)
	writeCompressedObject(t, "base", []byte("original\n"))
	writeCompressedObject(t, "broken", []byte("1\nnot a delta\n"))

	target := filepath.Join(t.TempDir(), "config")
	if err := os.WriteFile(target, []byte("working content\n"), 0o640); err != nil {
		t.Fatal(err)
	}

	val := tr.Tracker{
		Base:     "base",
		Versions: []tr.VersionDetails{{UID: "broken"}},
	}
	err := Reconstruct(val, target, LastVersion)
	if err == nil {
		t.Fatal("expected corrupt revision to fail")
	}

	content, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "working content\n" {
		t.Fatalf("target was modified after failed reconstruction: %q", content)
	}
	info, err := os.Stat(target)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o640 {
		t.Fatalf("target permissions changed: %o", info.Mode().Perm())
	}
}

func TestReconstructWritesRequestedRevisionAndPreservesMode(t *testing.T) {
	setupRepository(t)
	writeCompressedObject(t, "base", []byte("one\ntwo\n"))
	writeCompressedObject(t, "second", []byte("2\n2 @@@ dGhyZWU=\n"))

	target := filepath.Join(t.TempDir(), "config")
	if err := os.WriteFile(target, []byte("working\n"), 0o640); err != nil {
		t.Fatal(err)
	}

	val := tr.Tracker{
		Base:     "base",
		Versions: []tr.VersionDetails{{UID: "second"}},
	}
	if err := Reconstruct(val, target, LastVersion); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "one\nthree\n" {
		t.Fatalf("unexpected reconstructed content: %q", content)
	}
	info, err := os.Stat(target)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o640 {
		t.Fatalf("target permissions changed: %o", info.Mode().Perm())
	}
}

func TestReconstructRejectsSymlinkTarget(t *testing.T) {
	setupRepository(t)
	writeCompressedObject(t, "base", []byte("original\n"))

	dir := t.TempDir()
	actual := filepath.Join(dir, "actual")
	link := filepath.Join(dir, "link")
	if err := os.WriteFile(actual, []byte("keep\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(actual, link); err != nil {
		t.Fatal(err)
	}

	err := Reconstruct(tr.Tracker{Base: "base"}, link, BaseVersion)
	if err == nil {
		t.Fatal("expected symlink target to be rejected")
	}
	if errors.Is(err, os.ErrNotExist) {
		t.Fatalf("unexpected missing-file error: %v", err)
	}

	content, err := os.ReadFile(actual)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "keep\n" {
		t.Fatalf("symlink target was modified: %q", content)
	}
}
