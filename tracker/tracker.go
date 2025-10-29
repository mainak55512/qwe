package tracker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	cp "github.com/mainak55512/qwe/compressor"
	er "github.com/mainak55512/qwe/qwerror"
	utl "github.com/mainak55512/qwe/qweutils"
)

type VersionDetails struct {
	UID           string `json:"uid"`
	CommitMessage string `json:"commit_message"`
	TimeStamp     string `json:"time_stamp"`
}

type Tracker struct {
	Base     string           `json:"base"`
	Current  string           `json:"current"`
	Versions []VersionDetails `json:"versions"`
}

type FileDetails struct {
	FileName     string `json:"file_name"`
	CommitNumber int    `json:"commit_number"`
	FileObjID    string `json:"file_obj_id"`
}

type GroupVersionDetails struct {
	CommitMessage string                 `json:"commit_message"`
	Files         map[string]FileDetails `json:"files"`
}

type GroupTracker struct {
	GroupName    string                         `json:"group_name"`
	Current      string                         `json:"current"`
	VersionOrder []string                       `json:"version_order"`
	Versions     map[string]GroupVersionDetails `json:"versions"`
}

type TrackerSchema map[string]Tracker
type GroupTrackerSchema map[string]GroupTracker

// Returns the tracker details from _tracker.qwe or _group_tracker.qwe
func GetTracker(trackerType int) (TrackerSchema, GroupTrackerSchema, error) {
	var tracker_schema TrackerSchema
	var group_tracker_schema GroupTrackerSchema

	var trackerPath string

	// 0 is associated with file tracker, 1 is associated with group tracker
	if trackerType == 0 {
		trackerPath = ".qwe/_tracker.qwe"
	} else if trackerType == 1 {
		trackerPath = ".qwe/_group_tracker.qwe"
	} else {
		return nil, nil, er.InvalidTracker
	}

	// Decompress _tracker.qwe
	if err := cp.DecompressFile(trackerPath); err != nil {
		return nil, nil, err
	}

	file, err := os.Open(trackerPath)
	if err != nil {
		return nil, nil, err
	}

	reader := bufio.NewReader(file)
	current_tracker, err := io.ReadAll(reader)
	if err != nil {
		file.Close()
		return nil, nil, er.TrackerAccessErr
	} else {

		if trackerType == 0 {
			// Parse the content of the tracker file
			if err := json.Unmarshal(current_tracker, &tracker_schema); err != nil {
				file.Close()
				cp.CompressFile(trackerPath)
				return nil, nil, er.TrackerParseErr
			}
		} else {
			// Parse the content of the tracker file
			if err := json.Unmarshal(current_tracker, &group_tracker_schema); err != nil {
				file.Close()
				cp.CompressFile(trackerPath)
				return nil, nil, er.TrackerParseErr
			}
		}
	}
	file.Close()

	// Compress the tracker
	if err = cp.CompressFile(trackerPath); err != nil {
		return nil, nil, err
	}
	return tracker_schema, group_tracker_schema, nil
}

// Updates _tracker.qwe file
func SaveTracker(trackerType int, content []byte) error {

	var trackerPath string

	// 0 is associated with file tracker, 1 is associated with group tracker
	if trackerType == 0 {
		trackerPath = ".qwe/_tracker.qwe"
	} else if trackerType == 1 {
		trackerPath = ".qwe/_group_tracker.qwe"
	} else {
		return er.InvalidTracker
	}

	// Truncate the tracker file
	tracker_content, err := os.Create(trackerPath)
	if err != nil {
		return err
	}

	// Write the new content on the tracker file
	writer := bufio.NewWriter(tracker_content)
	_, err = writer.Write(content)
	if err != nil {
		return er.BaseWriteErr
	}
	if err = writer.Flush(); err != nil {
		return er.TrackerWriteErr
	}
	tracker_content.Close()

	// Compress the tracker file
	if err = cp.CompressFile(trackerPath); err != nil {
		return err
	}
	return nil
}

// Creates an entry for the file in Tracker and generates a base varient of the file
func StartTracking(filePath string) (string, error) {

	// Get tracker details
	tracker, _, err := GetTracker(0)
	if err != nil {
		return "", err
	}

	fileId := utl.Hasher(filePath)

	// This will be used as the name of the base file
	fileObjectId := "_base_" + utl.Hasher(fmt.Sprintf("%s%d", filePath, time.Now().UnixNano()))

	// If the file is already tracked then return error
	if _, ok := tracker[fileId]; ok {
		return "", er.FileTracked
	}

	base_content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("File not found: %s", filePath)
	}

	// (Need to change to a buffered writer) write the content of the file to the base varient
	if err := os.WriteFile(".qwe/_object/"+fileObjectId, base_content, 0644); err != nil {
		return "", er.TrackUnsuccessful
	}

	// Compress the base file
	if err = cp.CompressFile(".qwe/_object/" + fileObjectId); err != nil {
		return "", err
	}

	// Add tracker entry for the file
	tracker[fileId] = Tracker{
		Base:     fileObjectId,
		Current:  fileObjectId,
		Versions: []VersionDetails{},
	}

	marshalContent, err := json.MarshalIndent(tracker, "", " ")
	if err != nil {
		return "", er.CommitUnsuccessful
	}

	// Update the tracker
	if err = SaveTracker(0, marshalContent); err != nil {
		return "", err
	}
	fmt.Println("Started tracking", filePath)
	return fileObjectId, nil
}

// Start tracking a file in a group
func StartGroupTracking(groupName, filePath string) error {

	// Get tracker details
	_, groupTracker, err := GetTracker(1)
	if err != nil {
		return err
	}
	if utl.FolderExists(filePath) {
		err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && path != filePath {
				return filepath.SkipDir
			}
			if !info.IsDir() {
				groupTracker, err = fileTracker(path, groupName, groupTracker)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	} else {
		groupTracker, err = fileTracker(filePath, groupName, groupTracker)
		if err != nil {
			return err
		}
	}

	marshalContent, err := json.MarshalIndent(groupTracker, "", " ")
	if err != nil {
		return er.CommitUnsuccessful
	}

	// Update the tracker
	if err = SaveTracker(1, marshalContent); err != nil {
		return err
	}
	return nil
}

func fileTracker(filePath string, groupName string, groupTracker GroupTrackerSchema) (GroupTrackerSchema, error) {
	// Get tracker details
	tracker, _, err := GetTracker(0)
	if err != nil {
		return groupTracker, err
	}
	fileId := utl.Hasher(filePath)
	groupId := utl.Hasher(groupName)
	f, ok := tracker[fileId]
	if ok { // If the file is already tracked, get the current version and update the group tracker
		val, ok := groupTracker[groupId]
		if !ok {
			return groupTracker, er.InvalidGroup
		}
		_, ok = val.Versions[val.Current].Files[fileId]
		if ok {
			return groupTracker, fmt.Errorf("File %s is already tracked in group %s", filePath, groupName)
		}
		var commitNumber int

		// if current version of the file is a base file
		if strings.HasPrefix(f.Current, "_base_") {
			commitNumber = -2 // means for revert we need to revert back to base version
		} else {
			for i, elem := range f.Versions {
				if elem.UID == f.Current {
					commitNumber = i
					break
				}
			}
		}
		val.Versions[val.Current].Files[fileId] = FileDetails{
			FileName:     filePath,
			CommitNumber: commitNumber,
			FileObjID:    f.Current,
		}
		groupTracker[groupId] = val
	} else { // If file is not tracked, then track the file first
		fileObjectId, err := StartTracking(filePath)
		if err != nil {
			return groupTracker, err
		}
		val, ok := groupTracker[groupId]
		if !ok {
			return groupTracker, er.InvalidGroup
		}

		// As the file is first time tracked, the commit id is set to -2, that indicates, in case of revert, need to revert back to base version
		val.Versions[val.Current].Files[fileId] = FileDetails{
			FileName:     filePath,
			CommitNumber: -2,
			FileObjID:    fileObjectId,
		}
		groupTracker[groupId] = val
	}
	return groupTracker, nil
}
