package astparser

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	PhpExecutablePath string // Path to the PHP executable
	parserScript      string // Path to the extracted parse.php script
	tempParserDir     string // Temporary directory where parser files are extracted
)

// Init initializes the AST parser.
func Init(phpExePath string, fsys fs.FS) error {
	PhpExecutablePath = phpExePath

	// Extract embedded parser files to a temporary directory
	var err error
	tempParserDir, err = extractFiles(fsys)
	if err != nil {
		return fmt.Errorf("failed to extract parser files: %w", err)
	}

	// parse.php should now be at the root of the extracted tempDir
	parserScript = filepath.Join(tempParserDir, "parse.php")
	
	// Verify extracted parser script exists
	if _, err := os.Stat(parserScript); os.IsNotExist(err) {
		return fmt.Errorf("error: Extracted parse.php not found at %s: %v", parserScript, err)
	}
	return nil
}

// Cleanup removes the temporary parser directory.
func Cleanup() {
	if tempParserDir != "" {
		os.RemoveAll(tempParserDir)
		log.Printf("Cleaned up temporary parser directory: %s", tempParserDir)
	}
}

// extractFiles extracts all files from the provided filesystem to a temporary directory.
func extractFiles(fsys fs.FS) (string, error) {
	tempDir, err := os.MkdirTemp("", "ci4-parser-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == "." {
			return nil
		}

		targetPath := filepath.Join(tempDir, path)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		srcFile, err := fsys.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file %s from filesystem: %w", path, err)
		}
		defer srcFile.Close()

		dstFile, err := os.Create(targetPath)
		if err != nil {
			return fmt.Errorf("failed to create target file %s: %w", targetPath, err)
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return fmt.Errorf("failed to copy file %s to %s: %w", path, targetPath, err)
		}
		return nil
	})

	if err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to walk and extract files: %w", err)
	}

	log.Printf("Parser files extracted to: %s", tempDir)
	return tempDir, nil
}


// ExecutePhpParser runs the PHP parser script on a given file and returns its AST as JSON.
func ExecutePhpParser(filePath string) ([]byte, error) {
	cmd := exec.Command(PhpExecutablePath, parserScript, filePath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Improved: include both stdout and stderr in the error message
		return nil, fmt.Errorf("failed to run PHP parser for %s: %v. Stdout: %s. Stderr: %s", 
			filePath, err, strings.TrimSpace(out.String()), strings.TrimSpace(stderr.String()))
	}

	output := out.Bytes()
	if strings.Contains(string(output), "Parse error") {
        return nil, fmt.Errorf("PHP parser reported error for %s: %s", filePath, strings.TrimSpace(string(output)))
    }

	return output, nil
}
