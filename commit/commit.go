package commit

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"strings"
	tw "text/tabwriter"

	cp "github.com/mainak55512/qwe/compressor"
	er "github.com/mainak55512/qwe/qwerror"
	utl "github.com/mainak55512/qwe/qweutils"
	res "github.com/mainak55512/qwe/reconstruct"
	tr "github.com/mainak55512/qwe/tracker"
)

// Tracks the difference of the uncommitted file
func CommitUnit(filePath, message string) (string, int, error) {

	// Get tracking details from _tracker.qwe
	tracker, _, err := tr.GetTracker(0)
	if err != nil {
		return "", -3, err
	}

	// Create hash of file name, it will be used later to retrive file details from tracker
	fileId := utl.Hasher(filePath)

	// hash from file name and current time, will be used later as the file name of the commit
	fileObjectId := utl.Hasher(fmt.Sprintf("%s%d", filePath, time.Now().UnixNano()))

	var commitID int

	// Check if file is tracked
	if val, ok := tracker[fileId]; ok {
		target := ".qwe/_object/" + fileObjectId

		// This is the latest version of uncommitted file changes
		new_file, err := os.Open(filePath)
		if err != nil {
			return "", -3, err // -3 means unsuccessful
		}
		defer new_file.Close()

		// Reconstruct the file to the latest committed version
		// by applying all the changes to the base version
		if err = res.Reconstruct(val, target, -1); err != nil {
			return "", -3, err // -3 means unsuccessful
		}

		current_file, err := os.Open(target)
		if err != nil {
			return "", -3, err
		}

		current_scanner := bufio.NewScanner(current_file)

		new_scanner := bufio.NewScanner(new_file)

		var diff_content string

		// Find the difference between latest uncommitted and committed versions and store that in diff_content
		// difference is stored as <line-number> @@@ <new string value>
		line := 0
		for new_scanner.Scan() {
			line++
			current_scanner.Scan()
			if !bytes.Equal(current_scanner.Bytes(), new_scanner.Bytes()) {
				diff_content += fmt.Sprintf("%d @@@ %s\n", line, utl.ConvStrEnc(new_scanner.Text()))
			}
		}

		// This ensures no redundent commits are created for the file if there is no change
		if diff_content == "" {
			if !current_scanner.Scan() {
				os.Remove(target)
				return val.Versions[len(val.Versions)-1].UID, len(val.Versions) - 1, er.NoDiff
			}
		}

		// Adding total line number of uncommitted file on top of the diff_content
		// This line number will be used while reconstructing the file later
		diff_content = fmt.Sprintf("%d\n%s", line, diff_content)
		current_file.Close()

		output_content, err := os.Create(target)
		if err != nil {
			return "", -3, err // -3 means unsuccessful
		}

		output_writer := bufio.NewWriter(output_content)
		_, err = output_writer.WriteString(diff_content)
		if err != nil {
			return "", -3, er.BaseWriteErr // -3 means unsuccessful
		}
		if err = output_writer.Flush(); err != nil {
			return "", -3, er.OutputWriteErr // -3 means unsuccessful
		}
		output_content.Close()

		// Compressing the commit file
		if err = cp.CompressFile(target); err != nil {
			return "", -3, err // -3 means unsuccessful
		}

		// Update tracker
		val.Versions = append(val.Versions, tr.VersionDetails{
			UID:           fileObjectId,
			CommitMessage: message,
			TimeStamp:     time.Now().String()[:16],
		})
		val.Current = fileObjectId
		tracker[fileId] = val

		commitID = len(val.Versions) - 1

	} else {
		return "", -3, er.FileNotTracked // -3 means unsuccessful
	}

	// Save the updated tracker in _tracker.qwe
	marshalContent, err := json.MarshalIndent(tracker, "", " ")
	if err != nil {
		return "", -3, er.CommitUnsuccessful // -3 means unsuccessful
	}

	if err = tr.SaveTracker(0, marshalContent); err != nil {
		return "", -3, err // -3 means unsuccessful
	}

	fmt.Println("Committed", filePath, " successfully with commit id", commitID)
	return fileObjectId, commitID, nil
}

// Commit all file changes that are tracked in the group
func CommitGroup(groupName, commitMessage string) error {

	// Get group tracker
	_, groupTracker, err := tr.GetTracker(1)
	if err != nil {
		return err
	}

	groupID := utl.Hasher(groupName)
	groupObjID := utl.Hasher(fmt.Sprintf("%s%d", groupName, time.Now().UnixNano()))

	// Check if valid group
	gr, ok := groupTracker[groupID]
	if !ok {
		return er.InvalidGroup
	}

	// version order array maintains the order of commit history, appending new commit version here
	gr.VersionOrder = append(gr.VersionOrder, groupObjID)

	// Fetching the current group commit
	current, ok := gr.Versions[gr.Current]
	if !ok {
		return er.CurrentGrpErr
	}

	// newFiles contains the modified file details for the new commit
	newFiles := make(map[string]tr.FileDetails)

	for k := range current.Files {

		// Commit each and every file that is tracked in the group
		fileObjectID, commitID, err := CommitUnit(current.Files[k].FileName, commitMessage)

		// Do not treat it as error if there is no change in the file
		if err != nil && !errors.Is(err, er.NoDiff) {
			return err
		}

		// Add modified file details to newFiles
		newFiles[k] = tr.FileDetails{
			FileName:     current.Files[k].FileName,
			CommitNumber: commitID,
			FileObjID:    fileObjectID,
		}
	}

	// Update current version with the newly created commit in the group tracker
	gr.Current = groupObjID

	// Add new entry to the versions details of the group tracker
	gr.Versions[groupObjID] = tr.GroupVersionDetails{
		CommitMessage: commitMessage,
		Files:         newFiles,
	}

	commitID := len(gr.Versions) - 1

	// Update the group tracker with new details
	groupTracker[groupID] = gr

	marshalContent, err := json.MarshalIndent(groupTracker, "", " ")
	if err != nil {
		return er.CommitUnsuccessful
	}

	// Update the tracker
	if err = tr.SaveTracker(1, marshalContent); err != nil {
		return err
	}
	fmt.Println("Successfully committed to group", groupName, "with commit id", commitID)
	return nil
}

// Prints the commit history with CommitID, Commit message, and time stamp details
func GetCommitList(filePath string) error {

	// Get tracker details
	tracker, _, err := tr.GetTracker(0)
	if err != nil {
		return err
	}

	fileId := utl.Hasher(filePath)

	w := new(tw.Writer)
	w.Init(os.Stdout, 0, 0, 0, ' ', tw.TabIndent)

	// Loop through versions of the file and print commitID, commit message and time stamp for each entry
	for i, e := range tracker[fileId].Versions {
		fmt.Fprintln(w,
			fmt.Sprintf(
				"\nID:\t%d\nCommit Message:\t%s\nTime Stamp:\t%s\n", i, e.CommitMessage, e.TimeStamp,
			),
		)
		w.Flush()
	}
	return nil
}

// Shows list of all commits of the specified group
func GetGroupCommitList(groupName string) error {

	// Get group tracker
	_, groupTracker, err := tr.GetTracker(1)
	if err != nil {
		return err
	}

	groupID := utl.Hasher(groupName)

	// Check if valid group
	gr, ok := groupTracker[groupID]
	if !ok {
		return er.InvalidGroup
	}

	w := new(tw.Writer)
	w.Init(os.Stdout, 0, 0, 0, ' ', tw.TabIndent)

	// Print every version details
	i := 0
	for _, k := range gr.VersionOrder {
		fmt.Fprintln(w, fmt.Sprintf("\nID:\t%d\nCommit Message:\t%s\n", i, gr.Versions[k].CommitMessage))
		i++
	}
	w.Flush()
	return nil
}

// Shows the current checked out version of the file
func CurrentCommit(filePath string) error {

	// Get tracker details
	tracker, _, err := tr.GetTracker(0)
	if err != nil {
		return err
	}
	fileId := utl.Hasher(filePath)

	// Return error if the file is not tracked
	val, ok := tracker[fileId]
	if !ok {
		return er.FileNotTracked
	}

	// Get the current version of the file
	currentVersion := val.Current

	w := new(tw.Writer)
	w.Init(os.Stdout, 0, 0, 0, ' ', tw.TabIndent)

	// If the current checked out version is a base file, the print the base details
	if strings.HasPrefix(currentVersion, "_base_") {
		fmt.Fprintf(w, "\nCurrent Commit ID:\tbase\nCommit Message:\tBase version\n")
	} else {

		// Loop through the file versions, when current version is found print the details of commitID, commit message
		for i, e := range tracker[fileId].Versions {
			if e.UID == currentVersion {
				fmt.Fprintf(w, "\nCurrent Commit ID:\t%d\nCommit Message:\t%s\n", i, e.CommitMessage)
				break
			}
		}
	}
	w.Flush()
	return nil
}

// Prints current group commit details
func CurrentGroupCommit(groupName string) error {

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

	// Get the commit id of current version from the group tracker
	var commitID int
	for i, e := range val.VersionOrder {
		if e == val.Current {
			commitID = i
			break
		}
	}

	// Print current commit details
	w := new(tw.Writer)
	w.Init(os.Stdout, 0, 0, 0, ' ', tw.TabIndent)
	fmt.Fprintf(w, "\nName:\t %s\nCurrent Commit ID:\t %d\nCommit Message:\t %s\n", val.GroupName, commitID, val.Versions[val.Current].CommitMessage)
	files := val.Versions[val.Current].Files
	fmt.Fprintf(w, "\nAssociated files:\n")
	for e := range files {
		fmt.Fprintf(w, "File: %s, \tCommitID: %d\n", files[e].FileName, files[e].CommitNumber)
	}
	w.Flush()
	return nil
}

// Shows the names of groups tracked in the current repository
func GroupNameList() error {
	// Get group tracker
	_, groupTracker, err := tr.GetTracker(1)
	if err != nil {
		return err
	}

	// print the group names
	for k := range groupTracker {
		fmt.Println(groupTracker[k].GroupName)
	}
	return nil
}
