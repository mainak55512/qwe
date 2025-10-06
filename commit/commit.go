package commit

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	utl "github.com/mainak55512/qwe/qweutils"
	tr "github.com/mainak55512/qwe/tracker"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func CommitUnit(filePath, message string) error {
	tracker, err := tr.GetTracker()
	if err != nil {
		return fmt.Errorf("Can not retrieve Current version of %s", filePath)
	}
	fileId := utl.Hasher(filePath)
	fileObjectId := utl.Hasher(fmt.Sprintf("%s%d", filePath, time.Now().UnixNano()))

	if val, ok := tracker[fileId]; ok {
		target := ".qwe/_object/" + fileObjectId
		base_content, _ := os.ReadFile(".qwe/_object/" + val.Base)
		err = os.WriteFile(target, base_content, 0644)
		for _, elem := range val.Versions {
			diff_file, err := os.Open(".qwe/_object/" + elem.UID)
			if err != nil {
				log.Fatalf("Error opening file: %v", err)
			}
			defer diff_file.Close()
			diff_scanner := bufio.NewScanner(diff_file)

			base_file, err := os.Open(target)
			if err != nil {
				log.Fatalf("Error opening file: %v", err)
			}
			base_scanner := bufio.NewScanner(base_file)

			var output string
			diff_scanner.Scan()
			total_lines, _ := strconv.Atoi(diff_scanner.Text())
			diff_scanner.Scan()
			idx := 1
			for i := 0; i < total_lines; i++ {
				base_scanner.Scan()
				comp := strings.Split(diff_scanner.Text(), " @@@ ")
				line_number, _ := strconv.Atoi(comp[0])
				if idx == line_number {
					dec_str, err := utl.ConvStrDec(comp[1])
					if err != nil {
						return err
					}
					output += dec_str + "\n"
					diff_scanner.Scan()
				} else {
					output += base_scanner.Text() + "\n"
				}
				idx++
			}
			err = os.WriteFile(target, []byte(output), 0644)
			if err != nil {
				log.Fatal("Can not write to base file")
			}
			base_file.Close()
		}

		current_file, err := os.Open(target)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		current_scanner := bufio.NewScanner(current_file)

		new_file, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer new_file.Close()
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
		err = os.WriteFile(target, []byte(diff_content), 0644)
		if err != nil {
			return fmt.Errorf("Can not commit file!")
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

	if err = os.WriteFile(".qwe/_tracker.qwe", marshalContent, 0644); err != nil {
		return fmt.Errorf("Commit unsuccessful!")
	}

	return nil
}
