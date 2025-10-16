package cli

import (
	"fmt"
	"os"
	"strconv"
	tw "text/tabwriter"

	cm "github.com/mainak55512/qwe/commit"
	"github.com/mainak55512/qwe/diff"
	utl "github.com/mainak55512/qwe/qweutils"
	rb "github.com/mainak55512/qwe/rebase"
	rc "github.com/mainak55512/qwe/recover"
	rv "github.com/mainak55512/qwe/revert"
)

/*
Version details and available commands
*/
func helpText() {
	fmt.Println(
		`
       @@@@@@@@@@                                                                          
    @@@           @@@@@@@                                                                  
  @@                    @@@                                                               
@@                        @@                                                              
@@       @                @@@                                                              
@            @@@@        @@ @@@                                                            
@     @@@   @@  @     @@      @@                                                           
@    @  @@@ @ @@@      @@      @@                                                          
@     @@   @@ @        @        @@                                                         
@           @ @       @@         @                                                         
@@          @ @     @@@          @                                                         
 @@@        @@@   @@@            @                                                         
   @@@@      @@@@                                                                          
      @    @@@ @                                                                           
      @@  @@ @ @@@@                                                                        
      @@@@   @    @@                                                                       
       @@@@@ @@@@@@                                                                        
   _____        _______ 
  / _ \ \      / / ____|
 | | | \ \ /\ / /|  _|  
 | |_| |\ V  V / | |___ 
  \__\_\ \_/\_/  |_____|
		`,
	)
	w := new(tw.Writer)
	w.Init(os.Stdout, 0, 0, 0, ' ', tw.TabIndent)
	fmt.Println("Version: 0.1.6 - nightly")
	fmt.Println()
	fmt.Println("[COMMANDS]:")
	fmt.Fprintln(w, "qwe init\t[Initialize qwe in present directory]")
	fmt.Fprintln(w, "qwe track <file-path>\t[Start tracking a file]")
	fmt.Fprintln(w, "qwe list <file-path>\t[Get list of all commits on the file]")
	fmt.Fprintln(w, "qwe commit <file-path> \"<commit message>\"\t[Commit current version of the file to the version control]")
	fmt.Fprintln(w, "qwe revert <file-path> <commit-id>\t[Revert the file to a previous version]")
	fmt.Fprintln(w, "qwe current <file-path>\t[Get current commit details of the file]")
	fmt.Fprintln(w, "qwe recover <file-path>\t[Restore deleted file if earlier tracked]")
	fmt.Fprintln(w, "qwe rebase <file-path>\t[Revert back to base version of the file]")
	fmt.Fprintln(w, "qwe diff <file-path>\t[Shows difference between latest uncommitted version and latest committed version]")
	fmt.Fprintln(w, "qwe diff <file-path> <commit-id-1> <commit-id-2>\t[Shows difference between two commits]")
	fmt.Fprintln(w, "qwe diff <file-path> uncommitted <commit-id>\t[Shows difference between latest uncommitted version and commit-id version]")
	fmt.Fprintln(w)
	w.Flush()
}

/*
Handles command line arguments like init, track, commit, revert etc.
*/
func HandleArgs() error {
	command_list := os.Args[1:]

	if len(command_list) == 0 {
		helpText()
	} else {

		switch command_list[0] {
		case "init":
			{
				if len(command_list) != 1 {
					return fmt.Errorf("Init command doesn't take any argument")
				}
				if err := utl.Init(); err != nil {
					return err
				}
			}
		case "track":
			{
				if len(command_list) != 2 {
					return fmt.Errorf("Track command accepts one argument")
				}
				if err := utl.StartTracking(command_list[1]); err != nil {
					return err
				}
			}
		case "commit":
			{
				if len(command_list) != 3 {
					return fmt.Errorf("Commit command accepts two arguments")
				}
				if err := cm.CommitUnit(command_list[1], command_list[2]); err != nil {
					return err
				}
			}
		case "list":
			{
				if len(command_list) != 2 {
					return fmt.Errorf("List command accepts one argument")
				}
				if err := utl.GetCommitList(command_list[1]); err != nil {
					return err
				}
			}
		case "revert":
			{
				if len(command_list) != 3 {
					return fmt.Errorf("Revert command accepts two arguments")
				}
				commitNumber, err := strconv.Atoi(command_list[2])
				if err != nil {
					return fmt.Errorf("Not a valid commit number")
				}
				if err := rv.Revert(commitNumber, command_list[1]); err != nil {
					return err
				}
			}
		case "diff":
			{
				if len(command_list) != 2 && len(command_list) != 4 {
					return fmt.Errorf("diff command accepts one or three arguments")
				} else if len(command_list) == 4 {
					if err := diff.Diff(command_list[1], command_list[2], command_list[3]); err != nil {
						return err
					}
				} else {
					if err := diff.Diff(command_list[1], "", ""); err != nil {
						return err
					}
				}
			}
		case "current":
			{
				if len(command_list) != 2 {
					return fmt.Errorf("current command accepts one argument")
				}
				if err := utl.CurrentCommit(command_list[1]); err != nil {
					return err
				}
			}
		case "recover":
			{
				if len(command_list) != 2 {
					return fmt.Errorf("recover command accepts one argument")
				}
				if err := rc.Recover(command_list[1]); err != nil {
					return err
				}
			}
		case "rebase":
			{
				if len(command_list) != 2 {
					return fmt.Errorf("rebase command accepts one argument")
				}
				if err := rb.Rebase(command_list[1]); err != nil {
					return err
				}
			}
		default:
			{
				helpText()
			}
		}
	}
	return nil
}
