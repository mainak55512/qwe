package revert

import (
	"encoding/json"
	"fmt"

	utl "github.com/mainak55512/qwe/qweutils"
	rb "github.com/mainak55512/qwe/rebase"
	res "github.com/mainak55512/qwe/reconstruct"
	tr "github.com/mainak55512/qwe/tracker"
)

// Reverts the file to a specific version
func Revert(commitNumber int, filePath string) error {

	// Check if the file is present before reverting
	if exists := utl.FileExists(filePath); !exists {
		return fmt.Errorf("Invalid file path")
	}

	// Get tracker details
	tracker, _, err := tr.GetTracker(0)
	if err != nil {
		return fmt.Errorf("Can not retrieve Current version of %s", filePath)
	}
	fileId := utl.Hasher(filePath)

	// Check if the file is tracked
	if val, ok := tracker[fileId]; ok {
		target := filePath

		// Check if the commit number is valid
		if commitNumber < 0 || commitNumber > len(val.Versions) {
			return fmt.Errorf("Not a valid commit number")
		}

		// Reconstruct the file till the specific commit number
		if err = res.Reconstruct(val, target, commitNumber); err != nil {
			return err
		}

		// Update the current version of the file in tracker
		val.Current = val.Versions[commitNumber].UID
		tracker[fileId] = val
		marshalContent, err := json.MarshalIndent(tracker, "", " ")
		if err != nil {
			return fmt.Errorf("Commit unsuccessful!")
		}

		// Update the tracker
		if err = tr.SaveTracker(0, marshalContent); err != nil {
			return err
		}
	}
	fmt.Println("Successfully reverted", filePath, " back to commit", commitNumber)
	return nil
}

func RevertGroup(groupName string, commitID int) error {

	_, groupTracker, err := tr.GetTracker(1)
	if err != nil {
		return err
	}

	groupID := utl.Hasher(groupName)

	val, ok := groupTracker[groupID]
	if !ok {
		return fmt.Errorf("Invalid group!")
	}
	files := val.Versions[val.VersionOrder[commitID]].Files
	for k := range files {
		commitNumber := files[k].CommitNumber
		if commitNumber >= 0 {
			if err := Revert(commitNumber, files[k].FileName); err != nil {
				return err
			}
		} else if commitNumber == -2 {
			if err := rb.Rebase(files[k].FileName); err != nil {
				return err
			}
		}
	}
	val.Current = val.VersionOrder[commitID]
	groupTracker[groupID] = val
	marshalContent, err := json.MarshalIndent(groupTracker, "", " ")
	if err != nil {
		return fmt.Errorf("Commit unsuccessful!")
	}

	// Update the tracker
	if err = tr.SaveTracker(1, marshalContent); err != nil {
		return err
	}
	return nil
}
