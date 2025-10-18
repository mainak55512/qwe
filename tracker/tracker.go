package tracker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	cp "github.com/mainak55512/qwe/compressor"
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

type GroupVersionDetails struct {
	CommitMessage string            `json:"commit_message"`
	Files         map[string]string `json:"files"`
}

type GroupTracker struct {
	GroupName string                         `json:"group_name"`
	Current   string                         `json:"current"`
	Versions  map[string]GroupVersionDetails `json:"versions"`
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
