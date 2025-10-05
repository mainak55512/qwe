package qweutils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	tr "github.com/mainak55512/qwe/tracker"
	"io/fs"
	"os"
	"strings"
	"time"
)

func ConvStrEnc(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, " /@@@/ ", " %@@@% "), " @@@ ", " /@@@/ ")
}

func ConvStrDec(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, " /@@@/ ", " @@@ "), " %@@@%", " /@@@/ ")
}

func Hasher(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	hashByte := hasher.Sum(nil)
	return hex.EncodeToString(hashByte)[:32]
}

func GetCommitList(filePath string) error {
	tracker, err := tr.GetTracker()
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

func StartTracking(filePath string) error {
	tracker, err := tr.GetTracker()
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

	tracker[fileId] = tr.Tracker{
		Base:     fileObjectId,
		Current:  fileObjectId,
		Versions: []tr.VersionDetails{},
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

func exists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return true, nil
		}
		return false, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func Init() error {
	qwePath := ".qwe"
	if exists, _ := exists(qwePath); exists {
		return fmt.Errorf("Repository already initiated!")
	} else {
		if err := os.MkdirAll(qwePath+"/_object/", os.ModePerm); err != nil {
			return fmt.Errorf("Can not initiate repository!")
		}
		if _, err := os.Create(qwePath + "/_tracker.qwe"); err != nil {
			os.RemoveAll(qwePath)
			return fmt.Errorf("Can not initiate repository!")
		}
		os.WriteFile(qwePath+"/_tracker.qwe", []byte("{}"), 0644)
	}
	return nil
}
