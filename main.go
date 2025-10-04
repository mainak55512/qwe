package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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

func Hasher(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	hashByte := hasher.Sum(nil)
	return hex.EncodeToString(hashByte)[:32]
}

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

func GetCommitList(filePath string) error {
	tracker, err := GetTracker()
	if err != nil {
		return fmt.Errorf("Can not retrieve Current version of %s", filePath)
	}
	fileId := Hasher(filePath)
	for i, e := range tracker[fileId].Versions {
		fmt.Println(
			fmt.Sprintf(
				"\nID:\t%d\nCommit Message:\t%s\nTime Stamp:\t%s\n", i, e.CommitMessage, e.TimeStamp,
			),
		)
	}
	return nil
}

func Revert(commitNumber int, filePath string) error {
	tracker, err := GetTracker()
	if err != nil {
		return fmt.Errorf("Can not retrieve Current version of %s", filePath)
	}
	fileId := Hasher(filePath)
	fileObjectId := Hasher(fmt.Sprintf("%s%d", filePath, time.Now().UnixNano()))

	if val, ok := tracker[fileId]; ok {
		base_content, _ := os.ReadFile(".qwe/_object/" + val.Base)
		err = os.WriteFile(".qwe/_object/_base_"+fileObjectId, base_content, 0644)
		for i, elem := range val.Versions {
			if i > commitNumber {
				break
			}
			diff_file, err := os.Open(".qwe/_object/" + elem.UID)
			if err != nil {
				log.Fatalf("Error opening file: %v", err)
			}
			defer diff_file.Close()
			diff_scanner := bufio.NewScanner(diff_file)

			base_file, err := os.Open(".qwe/_object/_base_" + fileObjectId)
			if err != nil {
				log.Fatalf("Error opening file: %v", err)
			}
			base_scanner := bufio.NewScanner(base_file)

			var output string
			diff_scanner.Scan()
			total_lines, _ := strconv.Atoi(diff_scanner.Text())
			diff_scanner.Scan()
			// idx := 1
			for i := 0; i < total_lines; i++ {
				base_scanner.Scan()
				comp := strings.Split(diff_scanner.Text(), " @@@ ")
				line_number, _ := strconv.Atoi(comp[0])
				if i+1 == line_number {
					output += comp[1] + "\n"
					diff_scanner.Scan()
				} else {
					output += base_scanner.Text() + "\n"
				}
				// idx++
			}
			// fmt.Println(output)
			err = os.WriteFile(".qwe/_object/_base_"+fileObjectId, []byte(output), 0644)
			if err != nil {
				log.Fatal("Can not write to base file")
			}
			base_file.Close()
		}

		output_content, _ := os.ReadFile(".qwe/_object/_base_" + fileObjectId)
		err = os.WriteFile(filePath, output_content, 0644)
		os.Remove(".qwe/_object/_base_" + fileObjectId)

		val.Current = val.Versions[commitNumber].UID
		tracker[fileId] = val
		marshalContent, err := json.MarshalIndent(tracker, "", " ")
		if err != nil {
			return fmt.Errorf("Commit unsuccessful!")
		}

		if err = os.WriteFile(".qwe/_tracker.qwe", marshalContent, 0644); err != nil {
			return fmt.Errorf("Commit unsuccessful!")
		}
	}
	return nil
}

func CommitUnit(filePath, message string) error {
	tracker, err := GetTracker()
	if err != nil {
		return fmt.Errorf("Can not retrieve Current version of %s", filePath)
	}
	fileId := Hasher(filePath)
	fileObjectId := Hasher(fmt.Sprintf("%s%d", filePath, time.Now().UnixNano()))

	if val, ok := tracker[fileId]; ok {
		base_content, _ := os.ReadFile(".qwe/_object/" + val.Base)
		err = os.WriteFile(".qwe/_object/_base_"+fileObjectId, base_content, 0644)
		for _, elem := range val.Versions {
			diff_file, err := os.Open(".qwe/_object/" + elem.UID)
			if err != nil {
				log.Fatalf("Error opening file: %v", err)
			}
			defer diff_file.Close()
			diff_scanner := bufio.NewScanner(diff_file)

			base_file, err := os.Open(".qwe/_object/_base_" + fileObjectId)
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
					output += comp[1] + "\n"
					if diff_scanner.Scan() {
						comp = strings.Split(diff_scanner.Text(), " @@@ ")
						line_number, _ = strconv.Atoi(comp[0])
					}
				} else {
					output += base_scanner.Text() + "\n"
				}
				idx++
			}
			err = os.WriteFile(".qwe/_object/_base_"+fileObjectId, []byte(output), 0644)
			if err != nil {
				log.Fatal("Can not write to base file")
			}
			base_file.Close()
		}

		current_file, err := os.Open(".qwe/_object/_base_" + fileObjectId)
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
				diff_content += fmt.Sprintf("%d @@@ %s\n", line, new_scanner.Text())
			}
		}
		diff_content = fmt.Sprintf("%d\n%s", line, diff_content)
		current_file.Close()
		os.Remove(".qwe/_object/_base_" + fileObjectId)
		err = os.WriteFile(".qwe/_object/"+fileObjectId, []byte(diff_content), 0644)
		if err != nil {
			return fmt.Errorf("Can not commit file!")
		}

		val.Versions = append(val.Versions, VersionDetails{
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

func StartTracking(filePath string) error {
	tracker, err := GetTracker()
	if err != nil {
		return fmt.Errorf("Can not retrieve Current version of %s", filePath)
	}
	fileId := Hasher(filePath)
	fileObjectId := Hasher(fmt.Sprintf("%s%d", filePath, time.Now().UnixNano()))

	if _, ok := tracker[fileId]; ok {
		return fmt.Errorf("File is already being tracked")
	}

	base_content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("File not found: %s", filePath)
	}
	if err := os.WriteFile(".qwe/_object/"+fileObjectId, base_content, 0644); err != nil {
		return fmt.Errorf("Tracking unsuccessful!")
	}

	tracker[fileId] = Tracker{
		Base:     fileObjectId,
		Current:  fileObjectId,
		Versions: []VersionDetails{},
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

func main() {
	file_name := "test_new.json"
	// StartTracking(file_name)
	// GetCommitList(file_name)
	Revert(3, file_name)
	// if err := CommitUnit(file_name, "Fourth Commit"); err != nil {
	// 	fmt.Println(err)
	// }
}
