package qweutils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
	tw "text/tabwriter"
	"time"

	cp "github.com/mainak55512/qwe/compressor"
	tr "github.com/mainak55512/qwe/tracker"
)

// Encodes strings to base64
func ConvStrEnc(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Retrieves strings from base64
func ConvStrDec(str string) (string, error) {
	dec_str, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", fmt.Errorf("Failed to decode!")
	}
	return string(dec_str), nil
}

// creates Hash of a given string and returns first 32 characters as string
func Hasher(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	hashByte := hasher.Sum(nil)
	return hex.EncodeToString(hashByte)[:32]
}

// Prints the commit history with CommitID, Commit message, and time stamp details
func GetCommitList(filePath string) error {

	// Get tracker details
	tracker, _, err := tr.GetTracker(0)
	if err != nil {
		return fmt.Errorf("Can not retrieve Current version of %s", filePath)
	}

	fileId := Hasher(filePath)

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

// Creates an entry for the file in Tracker and generates a base varient of the file
func StartTracking(filePath string) error {

	// Get tracker details
	tracker, _, err := tr.GetTracker(0)
	if err != nil {
		return err
	}

	fileId := Hasher(filePath)

	// This will be used as the name of the base file
	fileObjectId := "_base_" + Hasher(fmt.Sprintf("%s%d", filePath, time.Now().UnixNano()))

	// If the file is already tracked then return error
	if _, ok := tracker[fileId]; ok {
		return fmt.Errorf("File is already being tracked")
	}

	base_content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("File not found: %s", filePath)
	}

	// (Need to change to a buffered writer) write the content of the file to the base varient
	if err := os.WriteFile(".qwe/_object/"+fileObjectId, base_content, 0644); err != nil {
		return fmt.Errorf("Tracking unsuccessful!")
	}

	// Compress the base file
	if err = cp.CompressFile(".qwe/_object/" + fileObjectId); err != nil {
		return err
	}

	// Add tracker entry for the file
	tracker[fileId] = tr.Tracker{
		Base:     fileObjectId,
		Current:  fileObjectId,
		Versions: []tr.VersionDetails{},
	}

	marshalContent, err := json.MarshalIndent(tracker, "", " ")
	if err != nil {
		return fmt.Errorf("Commit unsuccessful!")
	}

	// Update the tracker
	if err = tr.SaveTracker(0, marshalContent); err != nil {
		return err
	}
	fmt.Println("Started tracking", filePath)
	return nil
}

// Shows the current checked out version of the file
func CurrentCommit(filePath string) error {

	// Get tracker details
	tracker, _, err := tr.GetTracker(0)
	if err != nil {
		return err
	}
	fileId := Hasher(filePath)

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

// Checks if a folder exists
func FolderExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return true
		}
		return false
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return false
}

// Check if a file exists
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

// Initiates qwe repository
func Init() error {
	qwePath := ".qwe"

	// Check if qwe is already initialized
	if exists := FolderExists(qwePath); exists {
		return fmt.Errorf("Repository already initiated!")
	} else {

		// Create objects directory
		if err := os.MkdirAll(qwePath+"/_object/", os.ModePerm); err != nil {
			return fmt.Errorf("Can not initiate repository!")
		}
		// Create _tracker.qwe file
		if _, err := os.Create(qwePath + "/_tracker.qwe"); err != nil {
			os.RemoveAll(qwePath)
			return fmt.Errorf("Can not initiate repository!")
		}
		// Create _group_tracker.qwe file
		if _, err := os.Create(qwePath + "/_group_tracker.qwe"); err != nil {
			os.RemoveAll(qwePath)
			return fmt.Errorf("Can not initiate repository!")
		}
		// Initialize the tracker with '{}'
		if err := tr.SaveTracker(0, []byte("{}")); err != nil {
			return err
		}
		// Initialize the group tracker with '{}'
		if err := tr.SaveTracker(1, []byte("{}")); err != nil {
			return err
		}
	}
	fmt.Println("QWE initiated")
	return nil
}
