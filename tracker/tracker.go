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

type TrackerSchema map[string]Tracker

// Returns the tracker details from _tracker.qwe
func GetTracker() (TrackerSchema, error) {
	var tracker_schema TrackerSchema

	// Decompress _tracker.qwe
	if err := cp.DecompressFile(".qwe/_tracker.qwe"); err != nil {
		return nil, err
	}

	file, err := os.Open(".qwe/_tracker.qwe")
	if err != nil {
		return nil, fmt.Errorf("Can not open tracker!")
	}

	reader := bufio.NewReader(file)
	current_tracker, err := io.ReadAll(reader)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("Can not access tracker!")
	} else {

		// Parse the content of the tracker file
		if err := json.Unmarshal(current_tracker, &tracker_schema); err != nil {
			file.Close()
			cp.CompressFile(".qwe/_tracker.qwe")
			return nil, fmt.Errorf("Can not parse tracker!")
		}
	}
	file.Close()

	// Compress the tracker
	if err = cp.CompressFile(".qwe/_tracker.qwe"); err != nil {
		return nil, err
	}
	return tracker_schema, nil
}

func SaveTracker(content []byte) error {

	// Truncate the tracker file
	tracker_content, err := os.Create(".qwe/_tracker.qwe")
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
	if err = cp.CompressFile(".qwe/_tracker.qwe"); err != nil {
		return err
	}
	return nil
}
