package revert

import (
	// "bufio"
	"encoding/json"
	"fmt"
	// cp "github.com/mainak55512/qwe/compressor"
	utl "github.com/mainak55512/qwe/qweutils"
	tr "github.com/mainak55512/qwe/tracker"
	// "io"
	// "log"
	// "os"
	// "strconv"
	// "strings"
	res "github.com/mainak55512/qwe/reconstruct"
)

func Revert(commitNumber int, filePath string) error {
	tracker, err := tr.GetTracker()
	if err != nil {
		return fmt.Errorf("Can not retrieve Current version of %s", filePath)
	}
	fileId := utl.Hasher(filePath)

	if val, ok := tracker[fileId]; ok {
		target := filePath
		if commitNumber < 0 || commitNumber > len(val.Versions) {
			return fmt.Errorf("Not a valid commit number")
		}

		if err = res.Reconstruct(val, target, commitNumber); err != nil {
			return err
		}

		val.Current = val.Versions[commitNumber].UID
		tracker[fileId] = val
		marshalContent, err := json.MarshalIndent(tracker, "", " ")
		if err != nil {
			return fmt.Errorf("Commit unsuccessful!")
		}

		if err = tr.SaveTracker(marshalContent); err != nil {
			return err
		}
	}
	return nil
}
