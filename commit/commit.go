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

func CommitUnit(filePath, message string) error {
	tracker, err := tr.GetTracker()
	if err != nil {
		return fmt.Errorf("Can not retrieve Tracker, err: %s", err)
	}
	fileId := utl.Hasher(filePath)
	fileObjectId := utl.Hasher(fmt.Sprintf("%s%d", filePath, time.Now().UnixNano()))

	if val, ok := tracker[fileId]; ok {
		target := ".qwe/_object/" + fileObjectId

		new_file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("Error opening file: %v", err)
		}
		defer new_file.Close()

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

		line := 0
		for new_scanner.Scan() {
			line++
			current_scanner.Scan()
			if !bytes.Equal(current_scanner.Bytes(), new_scanner.Bytes()) {
				diff_content += fmt.Sprintf("%d @@@ %s\n", line, utl.ConvStrEnc(new_scanner.Text()))
			}
		}
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
		if err = cp.CompressFile(target); err != nil {
			return err
		}

		val.Versions = append(val.Versions, tr.VersionDetails{
			UID:           fileObjectId,
			CommitMessage: message,
			TimeStamp:     time.Now().String()[:16],
		})
		val.Current = fileObjectId
		tracker[fileId] = val

	} else {
		return fmt.Errorf("File is not tracked: %s", filePath)
	}

	marshalContent, err := json.MarshalIndent(tracker, "", " ")
	if err != nil {
		return fmt.Errorf("Commit unsuccessful!")
	}

	if err = tr.SaveTracker(marshalContent); err != nil {
		return err
	}

	fmt.Println("Committed successfully")
	return nil
}
