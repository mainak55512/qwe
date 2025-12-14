package tracker

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	cp "github.com/mainak55512/qwe/compressor"
	"github.com/mainak55512/qwe/qwerror"
	utl "github.com/mainak55512/qwe/qweutils"
)

const (
	FileName             = "_track_files.qwe"
	QweDir               = ".qwe"
	TrackFilePermissions = 0o644
)

// later add to configs and user might set itself
var excludedDirs = map[string]struct{}{
	".git": {},
	".qwe": {},
}

type TrackFiles map[string]*TrackFile

type TrackFile struct {
	FilePath string `json:"filepath"`
}

func NewTrackFile(filepath string) *TrackFile {
	return &TrackFile{
		FilePath: filepath,
	}
}

func (tf TrackFiles) Print() error {
	// Build a tree structure from file paths
	type node struct {
		name     string // name dir or file
		children map[string]*node
		isFile   bool
	}
	root := &node{
		name:     ".",
		children: make(map[string]*node),
		isFile:   false,
	}

	// insert each tracked filepath into the tree
	for _, t := range tf {
		// normalize to use '/' separators even on Windows earlier code uses '/'
		parts := strings.Split(strings.TrimPrefix(t.FilePath, "./"), "/")
		current := root
		for i, p := range parts {
			if p == "" {
				continue
			}
			child, childExists := current.children[p]
			if !childExists {
				child = &node{
					name:     p,
					children: make(map[string]*node),
				}
				current.children[p] = child
			}
			// mark leaf as file
			isLastPartOfPath := i == len(parts)-1
			if isLastPartOfPath {
				child.isFile = true
			}
			current = child
		}
	}

	// helper to get sorted keys
	sortKeys := func(m map[string]*node) []string {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return keys
	}

	// recursive print with prefix info for drawing lines
	var printNode func(n *node, prefix string, isLast bool)

	printNode = func(n *node, prefix string, isLast bool) {
		// don't print root name itself, print its children
		if n != root {
			var connector string
			if isLast {
				connector = "└── "
			} else {
				connector = "├── "
			}
			// colorize: directories bold, files blue
			reset := "\033[0m"
			bold := "\033[1m"
			blue := "\033[34m"
			if n.isFile {
				fmt.Printf("%s%s[Qwe] %s%s%s\n", prefix, connector, blue, n.name, reset)
			} else {
				fmt.Printf("%s%s%s%s%s\n", prefix, connector, bold, n.name, reset)
			}
		}
		// prepare next prefix
		var nextPrefix string
		if n == root {
			nextPrefix = ""
		} else {
			if isLast {
				nextPrefix = prefix + "    "
			} else {
				nextPrefix = prefix + "│   "
			}
		}

		keys := sortKeys(n.children)
		for i, k := range keys {
			child := n.children[k]
			last := i == len(keys)-1
			printNode(child, nextPrefix, last)
		}
	}

	// Print header and then children
	fmt.Println(".")
	keys := sortKeys(root.children)
	for i, k := range keys {
		child := root.children[k]
		last := i == len(keys)-1
		printNode(child, "", last)
	}
	return nil
}

func (tf TrackFiles) Save() error {
	bytes, err := json.MarshalIndent(tf, "", "  ")

	if err != nil {
		return fmt.Errorf("marshal tracked files: %w", err)
	}

	// atomic write
	filePath := filepath.Join(QweDir, FileName)
	tmp := filePath + ".tmp"

	if err := os.WriteFile(tmp, bytes, TrackFilePermissions); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}

	if err := os.Rename(tmp, filePath); err != nil {
		return fmt.Errorf("rename temp file: %w", err)
	}

	if err = cp.CompressFile(filePath); err != nil {
		return fmt.Errorf("compress file: %w", err)
	}
	return nil
}

func InitTrackedFiles() error {
	filePath := filepath.Join(QweDir, FileName)

	if err := os.WriteFile(filePath, []byte("{}"), TrackFilePermissions); err != nil {
		return fmt.Errorf("create tracked files file: %w", err)
	}

	if err := cp.CompressFile(filePath); err != nil {
		return fmt.Errorf("compress file: %w", err)
	}

	return nil
}

func UpdateTrackedFile(fileId, filePath string) error {
	trackFilePath := filepath.Join(QweDir, FileName)

	trackedFiles, err := LoadTrackedFilesFromFile(trackFilePath)

	if err != nil {
		return fmt.Errorf("error when try to loading tracked files: %w", err)
	}

	trackedFiles[fileId] = NewTrackFile(filePath)

	// Save the updated tracked files
	if err := trackedFiles.Save(); err != nil {
		return fmt.Errorf("save tracked files: %w", err)
	}

	return nil
}

func PrintTrackedFiles() error {
	if !utl.QweIsInWorkingDir() {
		return qwerror.RepoNotFound
	}

	trackFilePath := filepath.Join(QweDir, FileName)

	tf, err := loadOrScanTrackedFiles(trackFilePath)
	if err != nil {
		return err
	}

	return tf.Print()
}

func loadOrScanTrackedFiles(filePath string) (TrackFiles, error) {
	if utl.FileExists(filePath) {
		return LoadTrackedFilesFromFile(filePath)
	}
	workingDir, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	trackedFiles := make(TrackFiles)
	err = findTrackedFilesInCurrDir(workingDir, workingDir, &trackedFiles)

	if err != nil {
		return nil, err
	}
	err = trackedFiles.Save()

	if err != nil {
		return nil, fmt.Errorf("save tracked files: %w", err)
	}

	return trackedFiles, nil
}

func findTrackedFilesInCurrDir(dir, workingDir string, trackedFiles *TrackFiles) error {
	tracker, _, err := GetTracker(FileTrackerType)
	if err != nil {
		return err
	}

	fmt.Println("Scanning files in", dir, "...")

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				fmt.Printf("Permission denied: %s\n", path)
				return nil
			}
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}

		if d.IsDir() && !isExcludedDir(d.Name()) && path != dir {
			fmt.Println("Scanning files in", path, "...")
		}

		// Skip excluded directories
		if d.IsDir() && isExcludedDir(d.Name()) && path != dir {
			return fs.SkipDir
		}

		// Process files only
		if !d.IsDir() {
			relPath, err := filepath.Rel(workingDir, path)
			if err != nil {
				return err
			}
			fileId := utl.Hasher(relPath)
			if _, ok := tracker[fileId]; ok {
				(*trackedFiles)[fileId] = NewTrackFile(relPath)
			}
		}

		return nil
	})

	return err
}

// read data from file and load to memory (TrackFiles structure)
func LoadTrackedFilesFromFile(filename string) (TrackFiles, error) {
	trackedFiles := make(TrackFiles)

	if err := cp.DecompressFile(filename); err != nil {
		return nil, err
	}

	// If file doesn't exist, return empty map
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return trackedFiles, nil
	} else if err != nil {
		return nil, fmt.Errorf("stat file: %w", err)
	}

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	if err := json.Unmarshal(bytes, &trackedFiles); err != nil {
		return nil, fmt.Errorf("unmarshal tracked files: %w", err)
	}
	if err := cp.CompressFile(filename); err != nil {
		return nil, err
	}
	return trackedFiles, nil
}

func isExcludedDir(name string) bool {
	_, ok := excludedDirs[name]
	return ok
}
