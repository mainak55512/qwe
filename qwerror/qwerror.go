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
	CLIInitErr         = new(24, "init command doesn't take any argument!")
	CLIGrpInitErr      = new(25, "group-init command only takes 'group name' as argument!")
	CLITrackErr        = new(26, "track command only accepts 'file path' as argument!")
	CLIGrpTrackErr     = new(27, "group-track command accepts 'group name' and 'file path' as arguments!")
	CLICommitErr       = new(28, "commit command accepts 'file path' and 'commit message' as arguments!")
	CLIGrpCommitErr    = new(29, "group-commit command accepts 'group name' and 'commit message' as arguments!")
	CLIListErr         = new(30, "list command only accepts 'file path' as argument!")
	CLIGrpListErr      = new(31, "group-list command only accepts 'group name' as argument!")
	CLIRevertErr       = new(32, "Revert command either accepts no argument or two mandatory arguments 'group name' and 'commit number'!")
	CLIGrpRevertErr    = new(33, "group-revert command accepts 'group name' and 'commit number' as arguments!")
	CLIDiffErr         = new(34, "diff command accepts 'file path' as argument or 'file path' and two commit numbers as arguments!")
	CLICurrentErr      = new(35, "current command only accepts 'file path' as argument!")
	CLIGrpCurrentErr   = new(36, "group-current command only accepts 'group name' as argument!")
	CLIRecoverErr      = new(37, "recover command only accepts 'file path' as argument!")
	CLIRebaseErr       = new(38, "rebase command only accepts 'file path' as argument!")
	NoFileOrDiff       = new(39, "File does not exist or no changes found with the previous commit!")
	GrpNameListErr     = new(40, "groups command takes no arguments!")
)
