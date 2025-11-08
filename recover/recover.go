package recover

import (
	"fmt"
	"strings"

	bh "github.com/mainak55512/qwe/binaryhandler"
	er "github.com/mainak55512/qwe/qwerror"
	utl "github.com/mainak55512/qwe/qweutils"
	res "github.com/mainak55512/qwe/reconstruct"
	tr "github.com/mainak55512/qwe/tracker"
)

// Restores a deleted file if it was earlier tracked by qwe
func Recover(filePath string) error {

	// Check if the file is present before recovering
	if exists := utl.FileExists(filePath); exists {
		return er.FileExists
	}

	// Get tracker details
	tracker, _, err := tr.GetTracker(0)
	if err != nil {
		return err
	}
	fileId := utl.Hasher(filePath)
	val, ok := tracker[fileId]
	if !ok {
		return er.FileNotTracked
	}

	target := filePath

	// isBin, err := bh.CheckBinFile(filePath)
	// if err != nil {
	// 	return err
	// }

	if strings.HasPrefix(val.Base, "_bin_") {
		if err := bh.RevertBinFile(filePath, val.Current); err != nil {
			return err
		}
	} else {
		// Reconstruct the file all the way to the latest version
		if err = res.Reconstruct(val, target, -1); err != nil {
			return err
		}
	}

	fmt.Println("Successfully recovered", filePath)
	return nil
}
