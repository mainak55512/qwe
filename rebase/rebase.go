package rebase

import (
	"encoding/json"
	"fmt"
	er "github.com/mainak55512/qwe/qwerror"
	utl "github.com/mainak55512/qwe/qweutils"
	res "github.com/mainak55512/qwe/reconstruct"
	tr "github.com/mainak55512/qwe/tracker"
)

// Reverts a file back to its base version
func Rebase(filePath string) error {

	// Get tracker details
	tracker, _, err := tr.GetTracker(0)
	if err != nil {
		return err
	}

	fileId := utl.Hasher(filePath)

	// Check if file is tracked
	val, ok := tracker[fileId]
	if !ok {
		return er.FileNotTracked
	}

	// Reconstruct the file till its base version
	if err = res.Reconstruct(val, filePath, -2); err != nil {
		return err
	}

	// Update the current version of the file in tracker
	val.Current = val.Base
	tracker[fileId] = val
	marshalContent, err := json.MarshalIndent(tracker, "", " ")
	if err != nil {
		return er.CommitUnsuccessful
	}

	// Update the tracker
	if err = tr.SaveTracker(0, marshalContent); err != nil {
		return err
	}
	fmt.Println("Successfully reverted", filePath, "back to base version")
	return nil
}
