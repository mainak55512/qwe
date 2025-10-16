package revert

import (
	"encoding/json"
	"fmt"
	utl "github.com/mainak55512/qwe/qweutils"
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
	tracker, err := tr.GetTracker()
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
		if err = tr.SaveTracker(marshalContent); err != nil {
			return err
		}
	}
	fmt.Println("Successfully reverted back to commit", commitNumber)
	return nil
}
