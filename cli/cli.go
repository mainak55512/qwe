package cli

import (
	"fmt"
	"os"
	"strconv"
	tw "text/tabwriter"

	cm "github.com/mainak55512/qwe/commit"
	"github.com/mainak55512/qwe/diff"
	in "github.com/mainak55512/qwe/initializer"
	er "github.com/mainak55512/qwe/qwerror"
	rb "github.com/mainak55512/qwe/rebase"
	rc "github.com/mainak55512/qwe/recover"
	rv "github.com/mainak55512/qwe/revert"
	tr "github.com/mainak55512/qwe/tracker"
)

/*
Version details and available commands
*/
func helpText() {
	fmt.Println(`
                                                                                     
                    @@@@@@@@@                                                        
               @@@@@@@@@@@@@@@@@@@@                                                  
            @@@@@@              @@@@@@                                               
         @@@@@                      @@@@@@@@@@@@@                                    
        @@@@                           @@@@@@@@@@@@@                                 
      @@@@                                       @@@@                                
     @@@@                                  @@      @@@                               
    @@@                                    @@       @@                               
   @@@                                              @@@                              
   @@@                    @@@@@                    @@@@@@                            
  @@@                     @@  @@                  @@@@@@@@@                          
  @@@                    @@@@@@@            @@@@@@@@   @@@@@                         
  @@         @@@       @@@@                @@@@@@        @@@@@                       
  @@       @@   @     @@@        @@@@@     @@@            @@@@@                      
  @@       @@   @   @@@         @@   @@    @@@              @@@@                     
  @@@       @@@@@  @@@          @@@@@@    @@@                @@@                     
  @@@        @@@   @@@        @@@       @@@@                 @@@@                    
   @@@        @@@  @@@      @@@        @@@@                   @@@                    
    @@@        @@@@@@@    @@@@      @@@@@                     @@@                    
     @@@@        @@@@@   @@@     @@@@@@                       @@@                    
      @@@@@@       @@@  @@  @@@@@@@@                                                 
         @@@@@@@@  @@@  @@ @@@@@                                                     
             @@@@  @@@  @@                                                           
              @@@  @@@  @@                                                           
              @@@  @@@  @@          @@@                                              
              @@@  @@@  @@@@@@@@@@@@@ @@                                             
              @@@@@@     @@@@@@@@@@@@@@@                                             
              @@@@@                 @@@                                              
              @@@@@@@@@                                                              
                                                                                     
                                                                                     
          @                                     @@                                   
      @@@@@@@@@@@   @@       @@@       @@    @@@@@@@@@                               
    @@        @@@    @@      @@@@     @@   @@        @@                              
   @@          @@     @@    @@ @@     @@   @@         @@                             
   @@           @     @@   @@   @@   @@   @@@@@@@@@@@@@@                             
   @@          @@      @@ @@     @@ @@     @@                                        
    @@@       @@@       @@@@     @@ @@     @@@       @@                              
      @@@@@@@@@@@       @@@       @@@        @@@@@@@@@                               
	       @@                                                                    
               @@                                                                    
               @@                                                                    
                                                                                     
		`)
	w := new(tw.Writer)
	w.Init(os.Stdout, 0, 0, 0, ' ', tw.TabIndent)
	fmt.Println("Version: v0.2.1")
	fmt.Println()
	fmt.Println("[COMMANDS]:")
	fmt.Fprintln(w, "qwe init\t[Initialize qwe in present directory]")
	fmt.Fprintln(w, "qwe group-init <group name>\t[Initialize a group to track multiple files]")
	fmt.Fprintln(w, "qwe track <file-path>\t[Start tracking a file]")
	fmt.Fprintln(w, "qwe group-track <group name> <file-path>\t[Start tracking a file in a group]")
	fmt.Fprintln(w, "qwe list <file-path>\t[Get list of all commits on the file]")
	fmt.Fprintln(w, "qwe group-list <group name>\t[Get list of all commits on the group]")
	fmt.Fprintln(w, "qwe commit <file-path> \"<commit message>\"\t[Commit current version of the file to the version control]")
	fmt.Fprintln(w, "qwe group-commit <group name> \"<commit message>\"\t[Commit current version of all the files tracked in the group]")
	fmt.Fprintln(w, "qwe revert <file-path>\t[Revert the file to the last committed version]")
	fmt.Fprintln(w, "qwe revert <file-path> <commit-id>\t[Revert the file to a previous version]")
	fmt.Fprintln(w, "qwe group-revert <group name> <commit-id>\t[Revert all the files tracked in the group to a previous version]")
	fmt.Fprintln(w, "qwe current <file-path>\t[Get current commit details of the file]")
	fmt.Fprintln(w, "qwe group-current <group name>\t[Get current commit details of the group]")
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
					return er.CLIInitErr
				}
				if err := in.Init(); err != nil {
					return err
				}
			}
		case "group-init":
			{
				if len(command_list) != 2 {
					return er.CLIGrpInitErr
				}
				if err := in.GroupInit(command_list[1]); err != nil {
					return err
				}
			}
		case "track":
			{
				if len(command_list) != 2 {
					return er.CLITrackErr
				}
				if _, err := tr.StartTracking(command_list[1]); err != nil {
					return err
				}
			}
		case "group-track":
			{
				if len(command_list) != 3 {
					return er.CLIGrpTrackErr
				}
				if err := tr.StartGroupTracking(command_list[1], command_list[2]); err != nil {
					return err
				}
			}
		case "commit":
			{
				if len(command_list) != 3 {
					return er.CLICommitErr
				}
				if _, _, err := cm.CommitUnit(command_list[1], command_list[2]); err != nil {
					return err
				}
			}
		case "group-commit":
			{
				if len(command_list) != 3 {
					return er.CLIGrpCommitErr
				}
				if err := cm.CommitGroup(command_list[1], command_list[2]); err != nil {
					return err
				}
			}
		case "list":
			{
				if len(command_list) != 2 {
					return er.CLIListErr
				}
				if err := cm.GetCommitList(command_list[1]); err != nil {
					return err
				}
			}
		case "group-list":
			{
				if len(command_list) != 2 {
					return er.CLIGrpListErr
				}
				if err := cm.GetGroupCommitList(command_list[1]); err != nil {
					return err
				}
			}
		case "revert":
			{
				if len(command_list) != 3 && len(command_list) != 2 {
					return er.CLIRevertErr
				}
				var commitNumber int
				var err error
				if len(command_list) == 3 {
					commitNumber, err = strconv.Atoi(command_list[2])
					if err != nil {
						return er.InvalidCommitNo
					}
				} else {
					commitNumber = -1
				}
				if err := rv.Revert(commitNumber, command_list[1]); err != nil {
					return err
				}
			}
		case "group-revert":
			{
				if len(command_list) != 3 {
					return er.CLIGrpRevertErr
				}
				commitNumber, err := strconv.Atoi(command_list[2])
				if err != nil {
					return er.InvalidCommitNo
				}
				if err := rv.RevertGroup(command_list[1], commitNumber); err != nil {
					return err
				}
			}
		case "diff":
			{
				if len(command_list) != 2 && len(command_list) != 4 {
					return er.CLIDiffErr
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
					return er.CLICurrentErr
				}
				if err := cm.CurrentCommit(command_list[1]); err != nil {
					return err
				}
			}
		case "group-current":
			{
				if len(command_list) != 2 {
					return er.CLIGrpCurrentErr
				}
				if err := cm.CurrentGroupCommit(command_list[1]); err != nil {
					return err
				}
			}
		case "recover":
			{
				if len(command_list) != 2 {
					return er.CLIRecoverErr
				}
				if err := rc.Recover(command_list[1]); err != nil {
					return err
				}
			}
		case "rebase":
			{
				if len(command_list) != 2 {
					return er.CLIRebaseErr
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
