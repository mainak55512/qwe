package tracker

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	er "github.com/mainak55512/qwe/qwerror"
	utl "github.com/mainak55512/qwe/qweutils"
)

// StopTracking removes a file from individual and group tracking.
// The working file and stored objects are left untouched.
func StopTracking(filePath string) error {
	if !utl.QweIsInWorkingDir() {
		return er.RepoNotFound
	}

	tracker, _, err := GetTracker(FileTrackerType)
	if err != nil {
		return err
	}

	fileID := utl.Hasher(filePath)
	if _, ok := tracker[fileID]; !ok {
		return er.FileNotTracked
	}

	_, groupTracker, err := GetTracker(GroupTrackerType)
	if err != nil {
		return err
	}

	trackedFiles, err := loadOrScanTrackedFiles(filepath.Join(QweDir, FileName))
	if err != nil {
		return err
	}

	delete(tracker, fileID)
	delete(trackedFiles, fileID)

	for groupID, group := range groupTracker {
		for versionID, version := range group.Versions {
			delete(version.Files, fileID)
			group.Versions[versionID] = version
		}
		groupTracker[groupID] = group
	}

	trackerContent, err := json.MarshalIndent(tracker, "", " ")
	if err != nil {
		return er.CommitUnsuccessful
	}
	groupTrackerContent, err := json.MarshalIndent(groupTracker, "", " ")
	if err != nil {
		return er.CommitUnsuccessful
	}

	// Save the individual tracker last so a partial write can be retried.
	if err := SaveTracker(GroupTrackerType, groupTrackerContent); err != nil {
		return err
	}
	if err := trackedFiles.Save(); err != nil {
		return err
	}
	if err := SaveTracker(FileTrackerType, trackerContent); err != nil {
		return err
	}

	fmt.Println("Stopped tracking", filePath)
	return nil
}
