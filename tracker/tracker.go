package tracker

import (
	"encoding/json"
	"fmt"
	"os"
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
	if current_tracker, err := os.ReadFile(".qwe/_tracker.qwe"); err != nil {
		return nil, fmt.Errorf("Can not access tracker!")
	} else {
		if err := json.Unmarshal(current_tracker, &tracker_schema); err != nil {
			return nil, fmt.Errorf("Can not parse tracker!")
		}
	}
	return tracker_schema, nil
}
