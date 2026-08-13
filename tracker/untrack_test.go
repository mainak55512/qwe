package tracker

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	er "github.com/mainak55512/qwe/qwerror"
	utl "github.com/mainak55512/qwe/qweutils"
)

func TestStopTracking(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	if err := os.MkdirAll(filepath.Join(QweDir, "_object"), 0o755); err != nil {
		t.Fatalf("failed to create object directory: %v", err)
	}
	if err := SaveTracker(FileTrackerType, []byte("{}")); err != nil {
		t.Fatalf("failed to initialize file tracker: %v", err)
	}
	if err := SaveTracker(GroupTrackerType, []byte("{}")); err != nil {
		t.Fatalf("failed to initialize group tracker: %v", err)
	}
	if err := InitTrackedFiles(); err != nil {
		t.Fatalf("failed to initialize tracked files: %v", err)
	}

	targetPath := "notes.txt"
	keptPath := "keep.txt"
	if err := os.WriteFile(targetPath, []byte("keep this content"), 0o644); err != nil {
		t.Fatalf("failed to create target file: %v", err)
	}
	if err := os.WriteFile(keptPath, []byte("keep tracking this"), 0o644); err != nil {
		t.Fatalf("failed to create kept file: %v", err)
	}

	targetObjectID, err := StartTracking(targetPath)
	if err != nil {
		t.Fatalf("failed to track target file: %v", err)
	}
	keptObjectID, err := StartTracking(keptPath)
	if err != nil {
		t.Fatalf("failed to track kept file: %v", err)
	}

	targetID := utl.Hasher(targetPath)
	keptID := utl.Hasher(keptPath)
	firstGroupID := utl.Hasher("first")
	secondGroupID := utl.Hasher("second")
	groupTracker := GroupTrackerSchema{
		firstGroupID: {
			GroupName:    "first",
			Current:      "first-current",
			VersionOrder: []string{"first-base", "first-current"},
			Versions: map[string]GroupVersionDetails{
				"first-base": {
					Files: map[string]FileDetails{
						targetID: {FileName: targetPath, CommitNumber: -2, FileObjID: targetObjectID},
					},
				},
				"first-current": {
					Files: map[string]FileDetails{
						targetID: {FileName: targetPath, CommitNumber: -2, FileObjID: targetObjectID},
						keptID:   {FileName: keptPath, CommitNumber: -2, FileObjID: keptObjectID},
					},
				},
			},
		},
		secondGroupID: {
			GroupName:    "second",
			Current:      "second-current",
			VersionOrder: []string{"second-current"},
			Versions: map[string]GroupVersionDetails{
				"second-current": {
					Files: map[string]FileDetails{
						targetID: {FileName: targetPath, CommitNumber: -2, FileObjID: targetObjectID},
					},
				},
			},
		},
	}
	groupTrackerContent, err := json.MarshalIndent(groupTracker, "", " ")
	if err != nil {
		t.Fatalf("failed to marshal group tracker: %v", err)
	}
	if err := SaveTracker(GroupTrackerType, groupTrackerContent); err != nil {
		t.Fatalf("failed to save group tracker: %v", err)
	}

	if err := StopTracking(targetPath); err != nil {
		t.Fatalf("StopTracking() failed: %v", err)
	}

	tracker, _, err := GetTracker(FileTrackerType)
	if err != nil {
		t.Fatalf("failed to read file tracker: %v", err)
	}
	if _, ok := tracker[targetID]; ok {
		t.Error("target file remains in individual tracker")
	}
	if _, ok := tracker[keptID]; !ok {
		t.Error("unrelated file was removed from individual tracker")
	}

	_, updatedGroups, err := GetTracker(GroupTrackerType)
	if err != nil {
		t.Fatalf("failed to read group tracker: %v", err)
	}
	for _, group := range updatedGroups {
		for _, version := range group.Versions {
			if _, ok := version.Files[targetID]; ok {
				t.Errorf("target file remains in group %q", group.GroupName)
			}
		}
	}
	if _, ok := updatedGroups[firstGroupID].Versions["first-current"].Files[keptID]; !ok {
		t.Error("unrelated group membership was removed")
	}

	trackedFiles, err := LoadTrackedFilesFromFile(filepath.Join(QweDir, FileName))
	if err != nil {
		t.Fatalf("failed to read tracked-file index: %v", err)
	}
	if _, ok := trackedFiles[targetID]; ok {
		t.Error("target file remains in tracked-file index")
	}
	if _, ok := trackedFiles[keptID]; !ok {
		t.Error("unrelated file was removed from tracked-file index")
	}

	content, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("working file was removed: %v", err)
	}
	if string(content) != "keep this content" {
		t.Errorf("working file was modified: %q", content)
	}
	if _, err := os.Stat(filepath.Join(QweDir, "_object", targetObjectID)); err != nil {
		t.Errorf("stored object was removed: %v", err)
	}

	if err := StopTracking(targetPath); !errors.Is(err, er.FileNotTracked) {
		t.Errorf("expected FileNotTracked, got %v", err)
	}
}
