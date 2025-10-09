package diff

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mainak55512/qwe/qweutils"
	"github.com/mainak55512/qwe/tracker"
)

// Diff shows changes between the last committed and current version of a tracked file.
func Diff(filePath string) error {
	trk, err := tracker.GetTracker()
	if err != nil {
		return err
	}

	fileID := qweutils.Hasher(filePath)
	val, ok := trk[fileID]
	if !ok {
		return fmt.Errorf("file %s is not tracked", filePath)
	}
	if len(val.Versions) == 0 {
		return fmt.Errorf("no previous commits found for %s", filePath)
	}

	// --- Reconstruct last committed version into a temp file ---
	tmpFile, err := os.CreateTemp("", "qwe_diff_*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Copy base object
	baseFile, err := os.Open(".qwe/_object/" + val.Base)
	if err != nil {
		return fmt.Errorf("failed to open base object: %v", err)
	}
	if _, err := io.Copy(tmpFile, baseFile); err != nil {
		return fmt.Errorf("failed to copy base: %v", err)
	}
	baseFile.Close()

	// Apply each diff version sequentially
	for _, version := range val.Versions {
		diffFile, err := os.Open(".qwe/_object/" + version.UID)
		if err != nil {
			return fmt.Errorf("error opening diff file: %v", err)
		}

		// Read the current reconstructed file
		reconstructed, err := os.ReadFile(tmpFile.Name())
		if err != nil {
			diffFile.Close()
			return err
		}

		reconLines := strings.Split(string(reconstructed), "\n")
		diffScanner := bufio.NewScanner(diffFile)

		var output strings.Builder
		// Skip header lines (assuming first two are metadata)
		diffScanner.Scan()
		diffScanner.Scan()

		lineIdx := 0
		for diffScanner.Scan() && lineIdx < len(reconLines) {
			line := diffScanner.Text()
			if strings.Contains(line, "@@@") {
				comp := strings.Split(line, "@@@")
				decoded, _ := qweutils.ConvStrDec(comp[1])
				output.WriteString(decoded + "\n")
			} else {
				output.WriteString(reconLines[lineIdx] + "\n")
			}
			lineIdx++
		}

		for ; lineIdx < len(reconLines); lineIdx++ {
			output.WriteString(reconLines[lineIdx] + "\n")
		}

		diffFile.Close()
		os.WriteFile(tmpFile.Name(), []byte(output.String()), 0644)
	}

	// --- Compare reconstructed (last commit) vs current file ---
	lastCommitContent, _ := os.ReadFile(tmpFile.Name())
	currentContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read current file: %v", err)
	}

	lastLines := strings.Split(strings.TrimRight(string(lastCommitContent), "\n"), "\n")
	currLines := strings.Split(strings.TrimRight(string(currentContent), "\n"), "\n")

	maxOldWidth := 0
	for _, line := range lastLines {
		if len(line) > maxOldWidth {
			maxOldWidth = len(line)
		}
	}

	fmt.Printf("\n=== Diff View for %s ===\n\n", filePath)
	fmt.Printf("%-*s | %s\n", maxOldWidth, "Last Commit", "Current File")
	fmt.Printf("%s-+-%s\n", strings.Repeat("-", maxOldWidth), strings.Repeat("-", 50))

	maxLines := len(lastLines)
	if len(currLines) > maxLines {
		maxLines = len(currLines)
	}

	for i := 0; i < maxLines; i++ {
		var oldLine, newLine string
		if i < len(lastLines) {
			oldLine = lastLines[i]
		}
		if i < len(currLines) {
			newLine = currLines[i]
		}
		fmt.Printf("%-*s | %s\n", maxOldWidth, oldLine, newLine)
	}

	fmt.Printf("\n%s\n", strings.Repeat("-", maxOldWidth+55))
	fmt.Println("End of Diff")

	return nil
}
