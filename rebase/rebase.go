package rebase

import (
	"encoding/json"
	"fmt"
	utl "github.com/mainak55512/qwe/qweutils"
	res "github.com/mainak55512/qwe/reconstruct"
	tr "github.com/mainak55512/qwe/tracker"
)

func Rebase(filePath string) error {
	tracker, err := tr.GetTracker()
	if err != nil {
		return fmt.Errorf("Can not retrieve Tracker, err: %s", err)
	}
	fileId := utl.Hasher(filePath)

	val, ok := tracker[fileId]
	if !ok {
		return fmt.Errorf("File is not tracked!")
	}

	if err = res.Reconstruct(val, filePath, -2); err != nil {
		return err
	}
	val.Current = val.Base
	tracker[fileId] = val
	marshalContent, err := json.MarshalIndent(tracker, "", " ")
	if err != nil {
		return fmt.Errorf("Commit unsuccessful!")
	}

	if err = tr.SaveTracker(marshalContent); err != nil {
		return err
	}
	return nil
}
