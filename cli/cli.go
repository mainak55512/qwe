package cli

import (
	"fmt"
	"os"
	"strconv"
	tw "text/tabwriter"

	cm "github.com/mainak55512/qwe/commit"
	// "github.com/mainak55512/qwe/diff"
	utl "github.com/mainak55512/qwe/qweutils"
	rv "github.com/mainak55512/qwe/revert"
)

func helpText() {
	w := new(tw.Writer)
	w.Init(os.Stdout, 0, 0, 0, ' ', tw.TabIndent)
	fmt.Println("Version: 0.1.1")
	fmt.Println()
	fmt.Println("[COMMANDS]:")
	fmt.Fprintln(w, "qwe init\t[Initialize qwe in present directory]")
	fmt.Fprintln(w, "qwe track <file-path>\t[Start tracking a file]")
	fmt.Fprintln(w, "qwe list <file-path>\t[Get list of all commits on the file]")
	fmt.Fprintln(w, "qwe commit <file-path> \"<commit message>\"\t[Commit current version of the file to the version control]")
	fmt.Fprintln(w, "qwe revert <file-path> <commit-id>\t[Revert the file to a previous version]")
	fmt.Fprintln(w)
	w.Flush()
}

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
				// if len(command_list) != 2 {
				// 	return fmt.Errorf("diff command accepts filename as argument")
				// }
				// if err := diff.Diff(command_list[1]); err != nil {
				// 	return err
				// }
				helpText()
			}
		default:
			{
				helpText()
			}
		}
	}
	return nil
}
