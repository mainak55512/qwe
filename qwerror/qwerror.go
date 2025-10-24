package qwerror

import "fmt"

type qwError struct {
	code    int
	message string
}

func (q qwError) Error() string {
	return fmt.Sprintf("Error %d: %s", q.code, q.message)
}

func new(code int, message string) *qwError {
	return &qwError{
		code:    code,
		message: message,
	}
}

var (
	RepoAlreadyInit    = new(1, "Repository is already initiated!")
	RepoInitError      = new(2, "Can not initiate repository!")
	RepoNotFound       = new(3, "No qwe repository found!")
	GrpAlreadyTracked  = new(4, "Group is already being tracked!")
	CommitUnsuccessful = new(5, "Commit Unsuccessful!")
	InvalidTracker     = new(6, "Invalid Tracker Type!")
	TrackerAccessErr   = new(7, "Can not access tracker!")
	TrackerParseErr    = new(8, "Can not parse tracker!")
	BaseWriteErr       = new(9, "Can not write to base file!")
	TrackerWriteErr    = new(10, "Tracker file write error!")
	FileTracked        = new(11, "File is already being tracked!")
	TrackUnsuccessful  = new(12, "Tracking unsuccessful!")
	InvalidGroup       = new(13, "Invalid Group!")
	OutputWriteErr     = new(14, "Can not write to Output file!")
	FileNotTracked     = new(15, "File is not being tracked!")
	CurrentGrpErr      = new(16, "Can not retrieve current group version!")
	InvalidFile        = new(17, "Invalid file path!")
	InvalidCommitNo    = new(18, "Invalid commit number!")
	FileExists         = new(19, "File already exists!")
	CompOpenErr        = new(20, "Can not open file to compress!")
	CompBufInitErr     = new(21, "Can not initialize compression buffer!")
	BufCopyErr         = new(22, "Can not copy from/to compression buffer!")
	DecompBufInitErr   = new(23, "Can not initialize decompression buffer!")
)
