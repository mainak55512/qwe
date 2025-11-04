package revert

import (
	"encoding/json"
	"fmt"
	"strings"

	bh "github.com/mainak55512/qwe/binaryhandler"
	er "github.com/mainak55512/qwe/qwerror"
	utl "github.com/mainak55512/qwe/qweutils"
	rb "github.com/mainak55512/qwe/rebase"
	res "github.com/mainak55512/qwe/reconstruct"
	tr "github.com/mainak55512/qwe/tracker"
)

// Reverts the file to a specific version
func Revert(commitNumber int, filePath string) error {

	// Check if the file is present before reverting
	if exists := utl.FileExists(filePath); !exists {
		return fmt.Errorf("%w: %s\nUse 'recover' command to restore '%[2]s' if it was tracked earlier", er.InvalidFile, filePath)
	}

	// Get tracker details
	tracker, _, err := tr.GetTracker(0)
	if err != nil {
		return err
	}
	fileId := utl.Hasher(filePath)

	// Check if the file is tracked
	if val, ok := tracker[fileId]; ok {
		// Check if the commit number is valid
		if commitNumber < -1 || commitNumber > len(val.Versions)-1 {
			return er.InvalidCommitNo
		}

		if len(val.Versions) == 0 {
			return fmt.Errorf("File %s was never committed, use 'rebase' command to revert back to base version", filePath)
		}

		if strings.HasPrefix(val.Base, "_bin_") {
			commitID := commitNumber
			if commitNumber == -1 {
				commitID = len(val.Versions) - 1
			}

			fileObjID := val.Versions[commitID].UID

			if err = bh.RevertBinFile(filePath, fileObjID); err != nil {
				return err
			}
		} else {

			target := filePath

			// Reconstruct the file till the specific commit number
			if err = res.Reconstruct(val, target, commitNumber); err != nil {
				return err
			}
		}

		// Update the current version of the file in tracker
		// if commitID is -1 that means reverted back to latest commit
		if commitNumber == -1 {
			commitNumber = len(val.Versions) - 1
		}
		val.Current = val.Versions[commitNumber].UID
		tracker[fileId] = val
		marshalContent, err := json.MarshalIndent(tracker, "", " ")
		if err != nil {
			return er.CommitUnsuccessful
		}

		// Update the tracker
		if err = tr.SaveTracker(0, marshalContent); err != nil {
			return err
		}
	}
	fmt.Println("Successfully reverted", filePath, " back to commit", commitNumber)
	return nil
}

// Revert a group to any specific version
func RevertGroup(groupName string, commitID int) error {

	// Get group tracker
	_, groupTracker, err := tr.GetTracker(1)
	if err != nil {
		return err
	}

	groupID := utl.Hasher(groupName)

	// Check if valid group
	val, ok := groupTracker[groupID]
	if !ok {
		return er.InvalidGroup
	}

	if commitID < 0 || commitID > len(val.VersionOrder)-1 {
		return er.InvalidCommitNo
	}

	// Get all the file details of that specific version
	files := val.Versions[val.VersionOrder[commitID]].Files

	for k := range files {
		commitNumber := files[k].CommitNumber

		if commitNumber >= 0 { // commit number +ve means normal tracked file
			if err := Revert(commitNumber, files[k].FileName); err != nil {
				return err
			}
		} else if commitNumber == -2 { // commit number -2 means file is just tracked in qwe, no other commits are present, hence need to revert to base version
			if err := rb.Rebase(files[k].FileName); err != nil {
				return err
			}
		}
	}

	// Update current version with newly checked out version
	val.Current = val.VersionOrder[commitID]

	// Update group tracker with new values
	groupTracker[groupID] = val

	marshalContent, err := json.MarshalIndent(groupTracker, "", " ")
	if err != nil {
		return er.CommitUnsuccessful
	}

	// Update the tracker
	if err = tr.SaveTracker(1, marshalContent); err != nil {
		return err
	}
	return nil
}
