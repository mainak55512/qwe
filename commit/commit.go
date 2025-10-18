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
)

// Tracks the difference of the uncommitted file
func CommitUnit(filePath, message string) error {

	// Get tracking details from _tracker.qwe
	tracker, _, err := tr.GetTracker(0)
	if err != nil {
		return fmt.Errorf("Can not retrieve Tracker, err: %s", err)
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
			return fmt.Errorf("Error opening file: %v", err)
		}
		defer new_file.Close()

		// Reconstruct the file to the latest committed version
		// by applying all the changes to the base version
		if err = res.Reconstruct(val, target, -1); err != nil {
			return err
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
			return err
		}

		output_writer := bufio.NewWriter(output_content)
		_, err = output_writer.WriteString(diff_content)
		if err != nil {
			return fmt.Errorf("Can not write to base file")
		}
		if err = output_writer.Flush(); err != nil {
			return fmt.Errorf("Output file write error")
		}
		output_content.Close()

		// Compressing the commit file
		if err = cp.CompressFile(target); err != nil {
			return err
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
		return fmt.Errorf("File is not tracked: %s", filePath)
	}

	// Save the updated tracker in _tracker.qwe
	marshalContent, err := json.MarshalIndent(tracker, "", " ")
	if err != nil {
		return fmt.Errorf("Commit unsuccessful!")
	}

	if err = tr.SaveTracker(0, marshalContent); err != nil {
		return err
	}

	fmt.Println("Committed successfully with commit id", commitID)
	return nil
}
