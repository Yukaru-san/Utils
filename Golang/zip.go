package goutils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

/* Example

func main() {
	err := ZipFolder("{path}/test", "testzip.zip", []string{"test.exe"})
	if err != nil {
		fmt.Println(err.Error())
	}
	err = UnpackZip("testzip.zip", "./TestDir")
	if err != nil {
		fmt.Println(err.Error())
	}
}

*/

// UnpackZip unpacks the archive in the given path
func UnpackZip(archiveDir, targetDir string) error {

	var fileReader io.ReadCloser
	var targetFile *os.File

	// Create reader
	reader, err := zip.OpenReader(archiveDir)
	if err != nil {
		return err
	}

	// Create directories
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	// Loop and create files
	for _, file := range reader.File {
		path := filepath.Join(targetDir, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err = file.Open()
		if err != nil {
			return err
		}

		targetFile, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	// Close
	reader.Close()
	fileReader.Close()
	targetFile.Close()

	return nil
}

// ZipFolder searches the given dir and implements found files /// note that ignoredFiles includes directories!
func ZipFolder(sourcePath, zipName string, ignoredFiles []string) (string, error) {

	// Create and prepare zip
	zipPath := fmt.Sprint(sourcePath, string(filepath.Separator), zipName)
	zipfile, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}
	defer zipfile.Close()

	var ZipArchive = zip.NewWriter(zipfile)
	defer ZipArchive.Close()

	info, err := os.Stat(sourcePath)
	if err != nil {
		return "", nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(sourcePath)
	}

	fmt.Println("---starting to search files---")
	filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {

		// Check if file should be ignored
		if IsIgnored(zipName, ignoredFiles, path, info) {
			fmt.Println("   - ignoring", info.Name())
			return nil
		}
		fmt.Println("   + implementing", info.Name())

		if err != nil {
			return err
		}

		// Set file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Handle directory entries
		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, sourcePath))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		// Create file entry and fill it
		writer, err := ZipArchive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
	fmt.Println("---finished searching and packing---")

	return zipPath, err
}

// IsIgnored returns true if the file or folder should be ignored when packing
func IsIgnored(zipName string, ignoredFiles []string, path string, info os.FileInfo) bool {

	// Ignore the new zip itself
	if strings.Contains(path, zipName) {
		return true
	}

	// Ignore prefered files
	for _, i := range ignoredFiles {
		if strings.Contains(path, i) {
			return true
		}
	}

	// File shouldn't be ignored
	return false
}
