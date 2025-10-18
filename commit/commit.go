package commit

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	cp "github.com/mainak55512/qwe/compressor"
	utl "github.com/mainak55512/qwe/qweutils"
	res "github.com/mainak55512/qwe/reconstruct"
	tr "github.com/mainak55512/qwe/tracker"
	"strings"
	tw "text/tabwriter"
)

// Tracks the difference of the uncommitted file
func CommitUnit(filePath, message string) (string, error) {

	// Get tracking details from _tracker.qwe
	tracker, _, err := tr.GetTracker(0)
	if err != nil {
		return "", fmt.Errorf("Can not retrieve Tracker, err: %s", err)
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
			return "", fmt.Errorf("Error opening file: %v", err)
		}
		defer new_file.Close()

		// Reconstruct the file to the latest committed version
		// by applying all the changes to the base version
		if err = res.Reconstruct(val, target, -1); err != nil {
			return "", err
		}

		current_file, err := os.Open(target)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
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

		// Adding total line number of uncommitted file on top of the diff_content
		// This line number will be used while reconstructing the file later
		diff_content = fmt.Sprintf("%d\n%s", line, diff_content)
		current_file.Close()

		output_content, err := os.Create(target)
		if err != nil {
			return "", err
		}

		output_writer := bufio.NewWriter(output_content)
		_, err = output_writer.WriteString(diff_content)
		if err != nil {
			return "", fmt.Errorf("Can not write to base file")
		}
		if err = output_writer.Flush(); err != nil {
			return "", fmt.Errorf("Output file write error")
		}
		output_content.Close()

		// Compressing the commit file
		if err = cp.CompressFile(target); err != nil {
			return "", err
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
		return "", fmt.Errorf("File is not tracked: %s", filePath)
	}

	// Save the updated tracker in _tracker.qwe
	marshalContent, err := json.MarshalIndent(tracker, "", " ")
	if err != nil {
		return "", fmt.Errorf("Commit unsuccessful!")
	}

	if err = tr.SaveTracker(0, marshalContent); err != nil {
		return "", err
	}

	fmt.Println("Committed", filePath, " successfully with commit id", commitID)
	return fileObjectId, nil
}

func CommitGroup(groupName, commitMessage string) error {
	_, groupTracker, err := tr.GetTracker(1)
	if err != nil {
		return err
	}

	groupID := utl.Hasher(groupName)
	groupObjID := utl.Hasher(fmt.Sprintf("%s%d", groupName, time.Now().UnixNano()))
	gr, ok := groupTracker[groupID]
	if !ok {
		return fmt.Errorf("Invalid group!")
	}
	gr.VersionOrder = append(gr.VersionOrder, groupObjID)
	current, ok := gr.Versions[gr.Current]
	if !ok {
		return fmt.Errorf("Can not retrieve current group version!")
	}
	newFiles := make(map[string]tr.FileDetails)

	for k := range current.Files {
		fileObjectID, err := CommitUnit(current.Files[k].FileName, commitMessage)
		if err != nil {
			return err
		}
		newFiles[k] = tr.FileDetails{
			FileName:  current.Files[k].FileName,
			FileObjID: fileObjectID,
		}
	}

	gr.Current = groupObjID
	gr.Versions[groupObjID] = tr.GroupVersionDetails{
		CommitMessage: commitMessage,
		Files:         newFiles,
	}
	commitID := len(gr.Versions) - 1
	groupTracker[groupID] = gr

	marshalContent, err := json.MarshalIndent(groupTracker, "", " ")
	if err != nil {
		return fmt.Errorf("Commit unsuccessful!")
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
		return fmt.Errorf("Can not retrieve Current version of %s", filePath)
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

func GetGroupCommitList(groupName string) error {
	_, groupTracker, err := tr.GetTracker(1)
	if err != nil {
		return fmt.Errorf("Can not retrieve Current version of %s", groupName)
	}

	groupID := utl.Hasher(groupName)
	gr, ok := groupTracker[groupID]
	if !ok {
		return fmt.Errorf("Invalid group!")
	}

	w := new(tw.Writer)
	w.Init(os.Stdout, 0, 0, 0, ' ', tw.TabIndent)

	i := 0
	for _, k := range gr.VersionOrder {
		fmt.Fprintln(w, fmt.Sprintf("\nID:\t%d\nCommit Message:\t%s\n", i, gr.Versions[k].CommitMessage))
		i++
	}
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
		return fmt.Errorf("File is not tracked!")
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
