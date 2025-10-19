package tracker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	cp "github.com/mainak55512/qwe/compressor"
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

	if trackerType == 0 {
		trackerPath = ".qwe/_tracker.qwe"
	} else if trackerType == 1 {
		trackerPath = ".qwe/_group_tracker.qwe"
	} else {
		return nil, nil, fmt.Errorf("Invalid Tracker type!")
	}

	// Decompress _tracker.qwe
	if err := cp.DecompressFile(trackerPath); err != nil {
		return nil, nil, err
	}

	file, err := os.Open(trackerPath)
	if err != nil {
		return nil, nil, fmt.Errorf("Can not open tracker!")
	}

	reader := bufio.NewReader(file)
	current_tracker, err := io.ReadAll(reader)
	if err != nil {
		file.Close()
		return nil, nil, fmt.Errorf("Can not access tracker!")
	} else {

		if trackerType == 0 {
			// Parse the content of the tracker file
			if err := json.Unmarshal(current_tracker, &tracker_schema); err != nil {
				file.Close()
				cp.CompressFile(trackerPath)
				return nil, nil, fmt.Errorf("Can not parse tracker!")
			}
		} else {
			// Parse the content of the tracker file
			if err := json.Unmarshal(current_tracker, &group_tracker_schema); err != nil {
				file.Close()
				cp.CompressFile(trackerPath)
				return nil, nil, fmt.Errorf("Can not parse tracker!")
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
	if trackerType == 0 {
		trackerPath = ".qwe/_tracker.qwe"
	} else if trackerType == 1 {
		trackerPath = ".qwe/_group_tracker.qwe"
	} else {
		return fmt.Errorf("Invalid Tracker type!")
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
		return fmt.Errorf("Can not write to base file")
	}
	if err = writer.Flush(); err != nil {
		return fmt.Errorf("Tracker file write error")
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
		return "", fmt.Errorf("File is already being tracked")
	}

	base_content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("File not found: %s", filePath)
	}

	// (Need to change to a buffered writer) write the content of the file to the base varient
	if err := os.WriteFile(".qwe/_object/"+fileObjectId, base_content, 0644); err != nil {
		return "", fmt.Errorf("Tracking unsuccessful!")
	}

	// Compress the base file
	if err = cp.CompressFile(".qwe/_object/" + fileObjectId); err != nil {
		return fileObjectId, err
	}

	// Add tracker entry for the file
	tracker[fileId] = Tracker{
		Base:     fileObjectId,
		Current:  fileObjectId,
		Versions: []VersionDetails{},
	}

	marshalContent, err := json.MarshalIndent(tracker, "", " ")
	if err != nil {
		return fileObjectId, fmt.Errorf("Commit unsuccessful!")
	}

	// Update the tracker
	if err = SaveTracker(0, marshalContent); err != nil {
		return fileObjectId, err
	}
	fmt.Println("Started tracking", filePath)
	return fileObjectId, nil
}

func StartGroupTracking(groupName, filePath string) error {

	// Get tracker details
	tracker, _, err := GetTracker(0)
	if err != nil {
		return err
	}

	// Get tracker details
	_, groupTracker, err := GetTracker(1)
	if err != nil {
		return err
	}

	fileId := utl.Hasher(filePath)
	groupId := utl.Hasher(groupName)

	// If the file is already tracked then return error
	f, ok := tracker[fileId]
	if ok {
		val, ok := groupTracker[groupId]
		if !ok {
			return fmt.Errorf("Invalid group!")
		}
		_, ok = val.Versions[val.Current].Files[fileId]
		if ok {
			return fmt.Errorf("File %s is already tracked in group %s", filePath, groupName)
		}
		var commitNumber int
		if strings.HasPrefix(f.Current, "_base_") {
			commitNumber = -2
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
	} else {
		fileObjectId, err := StartTracking(filePath)
		if err != nil {
			return err
		}
		val, ok := groupTracker[groupId]
		if !ok {
			return fmt.Errorf("Invalid group!")
		}
		val.Versions[val.Current].Files[fileId] = FileDetails{
			FileName:     filePath,
			CommitNumber: -2,
			FileObjID:    fileObjectId,
		}
		groupTracker[groupId] = val
	}
	marshalContent, err := json.MarshalIndent(groupTracker, "", " ")
	if err != nil {
		return fmt.Errorf("Commit unsuccessful!")
	}

	// Update the tracker
	if err = SaveTracker(1, marshalContent); err != nil {
		return err
	}
	return nil
}
