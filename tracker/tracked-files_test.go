package tracker

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	cp "github.com/mainak55512/qwe/compressor"
)

// setupTestDir creates a temp directory with .qwe folder and changes to it.
func setupTestDir(t *testing.T) (tempDirPath string, cleanup func()) {
	t.Helper()

	// Save original working directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	// Create temp directory
	tempDirPath, err = os.MkdirTemp("", "qwe-tracker-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}

	// Change to temp directory
	if err := os.Chdir(tempDirPath); err != nil {
		os.RemoveAll(tempDirPath)
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	// Create .qwe directory
	qweDir := filepath.Join(tempDirPath, QweDir)
	if err := os.MkdirAll(qweDir, 0o755); err != nil {
		os.Chdir(originalDir)
		os.RemoveAll(tempDirPath)
		t.Fatalf("failed to create .qwe directory: %v", err)
	}

	// Return cleanup function
	cleanup = func() {
		os.Chdir(originalDir)
		os.RemoveAll(tempDirPath)
	}

	return tempDirPath, cleanup
}

// TestInitTrackedFiles_SuccessfulCreation tests that InitTrackedFiles creates file successfully
func TestInitTrackedFiles_SuccessfulCreation(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	if err := os.Chdir(tempDir); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	err := InitTrackedFiles()
	if err != nil {
		t.Fatalf("InitTrackedFiles() failed: %v", err)
	}

	// if err := os.Chdir(tempDir); err != nil {
	// 	os.RemoveAll(tempDir)
	// 	t.Fatalf("failed to change to temp directory: %v", err)
	// }
	//
	// Verify the file was created (compressed in place)
	filePath := filepath.Join(QweDir, FileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("file %s was not created", filePath)
	}

	// Decompress and verify content
	if err := cp.DecompressFile(filePath); err != nil {
		t.Fatalf("failed to decompress file: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expectedContent := "{}"
	if string(content) != expectedContent {
		t.Errorf("expected file content %q, got %q", expectedContent, string(content))
	}
}

// TestInitTrackedFiles_FilePermissions tests that InitTrackedFiles creates file with correct permissions
func TestInitTrackedFiles_FilePermissions(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	if err := os.Chdir(tempDir); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	// Save current umask and set it to 0 for this test
	// oldUmask := syscall.Umask(0)
	// defer syscall.Umask(oldUmask)
	err := InitTrackedFiles()
	if err != nil {
		t.Fatalf("InitTrackedFiles() failed: %v", err)
	}

	filePath := filepath.Join(QweDir, FileName)

	// Decompress to check permissions
	if err := cp.DecompressFile(filePath); err != nil {
		t.Fatalf("failed to decompress file: %v", err)
	}

	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	// expectedPerms := os.FileMode(TrackFilePermissions)
	// if info.Mode().Perm() != expectedPerms {
	// 	t.Errorf("expected permissions %v, got %v", expectedPerms, info.Mode().Perm())
	// }

	// changed the condition as system's umask may change permissions
	mode := info.Mode().Perm()
	if mode&0600 != 0600 {
		t.Errorf("owner should have read and write permissions, got %v", mode)
	}

	// Check that others don't have write
	if mode&0002 != 0 {
		t.Errorf("others should not have write permission, got %v", mode)
	}
}

// TestIsExcludedDir_ExcludedDirectories tests that isExcludedDir correctly identifies excluded directories
func TestIsExcludedDir_ExcludedDirectories(t *testing.T) {
	tests := []struct {
		name     string
		dirName  string
		expected bool
	}{
		{
			name:     "git directory",
			dirName:  ".git",
			expected: true,
		},
		{
			name:     "qwe directory",
			dirName:  ".qwe",
			expected: true,
		},
		{
			name:     "regular directory",
			dirName:  "src",
			expected: false,
		},
		{
			name:     "node_modules",
			dirName:  "node_modules",
			expected: false,
		},
		{
			name:     "empty string",
			dirName:  "",
			expected: false,
		},
		{
			name:     "hidden directory not excluded",
			dirName:  ".hidden",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isExcludedDir(tt.dirName)
			if result != tt.expected {
				t.Errorf("isExcludedDir(%q) = %v, expected %v", tt.dirName, result, tt.expected)
			}
		})
	}
}

// TestLoadTrackedFilesFromFile_WithData tests loading from a file with tracked files
func TestLoadTrackedFilesFromFile_WithData(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	if err := os.Chdir(tempDir); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	filePath := filepath.Join(QweDir, FileName)

	// Create test data
	testData := TrackFiles{
		"file1": &TrackFile{FilePath: "/path/to/file1.txt"},
		"file2": &TrackFile{FilePath: "/path/to/file2.go"},
		"file3": &TrackFile{FilePath: "/path/to/file3.md"},
	}

	// Write test data to file
	jsonData, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("failed to marshal test data: %v", err)
	}

	if err := os.WriteFile(filePath, jsonData, TrackFilePermissions); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Compress the file
	if err := cp.CompressFile(filePath); err != nil {
		t.Fatalf("failed to compress file: %v", err)
	}

	// Load tracked files
	trackedFiles, err := LoadTrackedFilesFromFile(filePath)
	if err != nil {
		t.Fatalf("LoadTrackedFilesFromFile() failed: %v", err)
	}

	if len(trackedFiles) != len(testData) {
		t.Errorf("expected %d tracked files, got %d", len(testData), len(trackedFiles))
	}

	// Verify each file
	for fileId, expectedFile := range testData {
		actualFile, exists := trackedFiles[fileId]
		if !exists {
			t.Errorf("file %s not found in loaded tracked files", fileId)
			continue
		}

		if actualFile.FilePath != expectedFile.FilePath {
			t.Errorf("file %s: expected path %s, got %s", fileId, expectedFile.FilePath, actualFile.FilePath)
		}
	}
}

// TestLoadTrackedFilesFromFile_NonExistentFile tests loading from a non-existent file
func TestLoadTrackedFilesFromFile_NonExistentFile(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	filePath := filepath.Join(QweDir, "non_existent.qwe")

	// LoadTrackedFilesFromFile tries to decompress first, which will fail for non-existent files
	_, err := LoadTrackedFilesFromFile(filePath)
	if err == nil {
		t.Fatal("expected error when loading non-existent file, got nil")
	}
}

// TestLoadTrackedFilesFromFile_InvalidJSON tests loading from a file with invalid JSON
func TestLoadTrackedFilesFromFile_InvalidJSON(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	if err := os.Chdir(tempDir); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	filePath := filepath.Join(QweDir, FileName)

	// Create file with invalid JSON
	if err := os.WriteFile(filePath, []byte("not valid json"), TrackFilePermissions); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Compress the file
	if err := cp.CompressFile(filePath); err != nil {
		t.Fatalf("failed to compress file: %v", err)
	}

	_, err := LoadTrackedFilesFromFile(filePath)
	if err == nil {
		t.Fatal("expected error when loading invalid JSON, got nil")
	}
}

// TestUpdateTrackedFile_NewFile tests adding a new tracked file
func TestUpdateTrackedFile_NewFile(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	if err := os.Chdir(tempDir); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	// Initialize tracked files first
	if err := InitTrackedFiles(); err != nil {
		t.Fatalf("InitTrackedFiles() failed: %v", err)
	}

	fileId := "test-file-id"
	filePath := "/path/to/test/file.txt"

	err := UpdateTrackedFile(fileId, filePath)
	if err != nil {
		t.Fatalf("UpdateTrackedFile() failed: %v", err)
	}

	// Load tracked files and verify
	trackFilePath := filepath.Join(QweDir, FileName)
	// trackFilePath := filepath.Join(tempDir, FileName)
	trackedFiles, err := LoadTrackedFilesFromFile(trackFilePath)
	if err != nil {
		t.Fatalf("failed to load tracked files: %v", err)
	}

	if len(trackedFiles) != 1 {
		t.Errorf("expected 1 tracked file, got %d", len(trackedFiles))
	}

	file, exists := trackedFiles[fileId]
	if !exists {
		t.Fatal("tracked file not found after update")
	}

	if file.FilePath != filePath {
		t.Errorf("expected file path %s, got %s", filePath, file.FilePath)
	}
}

// TestUpdateTrackedFile_UpdateExistingFile tests updating an existing tracked file
func TestUpdateTrackedFile_UpdateExistingFile(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	if err := os.Chdir(tempDir); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	// Initialize tracked files
	if err := InitTrackedFiles(); err != nil {
		t.Fatalf("InitTrackedFiles() failed: %v", err)
	}

	fileId := "test-file-id"
	originalPath := "/path/to/original/file.txt"
	updatedPath := "/path/to/updated/file.txt"

	// Add file first time
	if err := UpdateTrackedFile(fileId, originalPath); err != nil {
		t.Fatalf("first UpdateTrackedFile() failed: %v", err)
	}

	// Update the same file with new path
	if err := UpdateTrackedFile(fileId, updatedPath); err != nil {
		t.Fatalf("second UpdateTrackedFile() failed: %v", err)
	}

	// Load tracked files and verify
	trackFilePath := filepath.Join(tempDir, QweDir, FileName)
	trackedFiles, err := LoadTrackedFilesFromFile(trackFilePath)
	if err != nil {
		t.Fatalf("failed to load tracked files: %v", err)
	}

	if len(trackedFiles) != 1 {
		t.Errorf("expected 1 tracked file, got %d", len(trackedFiles))
	}

	file, exists := trackedFiles[fileId]
	if !exists {
		t.Fatal("tracked file not found after update")
	}

	if file.FilePath != updatedPath {
		t.Errorf("expected updated file path %s, got %s", updatedPath, file.FilePath)
	}
}

// TestUpdateTrackedFile_MultipleFiles tests adding multiple tracked files
func TestUpdateTrackedFile_MultipleFiles(t *testing.T) {
	tempDir, cleanup := setupTestDir(t)
	defer cleanup()

	if err := os.Chdir(tempDir); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("failed to change to temp directory: %v", err)
	}

	// Initialize tracked files
	if err := InitTrackedFiles(); err != nil {
		t.Fatalf("InitTrackedFiles() failed: %v", err)
	}

	files := map[string]string{
		"file1": "/path/to/file1.txt",
		"file2": "/path/to/file2.go",
		"file3": "/path/to/file3.md",
	}

	// Add multiple files
	for fileId, filePath := range files {
		if err := UpdateTrackedFile(fileId, filePath); err != nil {
			t.Fatalf("UpdateTrackedFile(%s) failed: %v", fileId, err)
		}
	}

	// Load tracked files and verify
	trackFilePath := filepath.Join(QweDir, FileName)
	trackedFiles, err := LoadTrackedFilesFromFile(trackFilePath)
	if err != nil {
		t.Fatalf("failed to load tracked files: %v", err)
	}

	if len(trackedFiles) != len(files) {
		t.Errorf("expected %d tracked files, got %d", len(files), len(trackedFiles))
	}

	// Verify each file
	for fileId, expectedPath := range files {
		file, exists := trackedFiles[fileId]
		if !exists {
			t.Errorf("file %s not found in tracked files", fileId)
			continue
		}

		if file.FilePath != expectedPath {
			t.Errorf("file %s: expected path %s, got %s", fileId, expectedPath, file.FilePath)
		}
	}
}
