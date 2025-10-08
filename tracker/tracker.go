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

func GetTracker() (TrackerSchema, error) {
	var tracker_schema TrackerSchema
	cp.DecompressFile(".qwe/_tracker.qwe")
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
		if err := json.Unmarshal(current_tracker, &tracker_schema); err != nil {
			file.Close()
			cp.CompressFile(".qwe/_tracker.qwe")
			return nil, fmt.Errorf("Can not parse tracker!")
		}
	}
	file.Close()
	cp.CompressFile(".qwe/_tracker.qwe")
	return tracker_schema, nil
}

func SaveTracker(content []byte) error {
	tracker_content, err := os.Create(".qwe/_tracker.qwe")
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(tracker_content)
	_, err = writer.Write(content)
	if err != nil {
		return fmt.Errorf("Can not write to base file")
	}
	if err = writer.Flush(); err != nil {
		return fmt.Errorf("Tracker file write error")
	}
	tracker_content.Close()
	cp.CompressFile(".qwe/_tracker.qwe")
	return nil
}
