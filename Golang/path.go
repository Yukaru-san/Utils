package goutils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

//GetHome returns the home directory of the current user
func GetHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err.Error())
		return ""
	}
	return home
}

// FindFilesInDirBySuffix finds all files within a directory ending with the given Suffix, if ignoreSpecific is given, those files will be ignored
func FindFilesInDirBySuffix(directoryPath string, suffix string, ignoreSpecific string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, suffix) {
			if len(ignoreSpecific) > 0 {
				if strings.HasSuffix(path, ignoreSpecific) {
					return nil
				}
			}
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

// GetLastPartOfPath returns a path's last part like a file's name
func GetLastPartOfPath(path string) string {
	return filepath.Base(path)
}

// GetNthPartOfPath returns the nth part from a path, going right to left (file to C:)
func GetNthPartOfPath(path string, n int) string {
	filePathSplit := strings.Split(path, string(filepath.Separator))
	return filePathSplit[len(filePathSplit)-n-1]
}

// ReplaceFilepathSeparator can be used to fix inconsistent separators
func ReplaceFilepathSeparator(filePath string, newSeparator string) string {
	filePath = strings.ReplaceAll(filePath, "\\\\", newSeparator)
	filePath = strings.ReplaceAll(filePath, "\\", newSeparator)
	filePath = strings.ReplaceAll(filePath, "/", newSeparator)
	if strings.HasSuffix(filePath, newSeparator) {
		filePath = filePath[:len(filePath)-1]
	}

	return filePath
}

// SanitizeOutput tries to create a correct output path across user inputs and devices
func SanitizeOutput(outputPath, fileName string) string {
	// Given path is a directory
	if !strings.Contains(GetLastPartOfPath(outputPath), ".") {
		// Create dir if needed
		os.MkdirAll(outputPath, 0750)

		// Empty string
		if len(fileName) == 0 {
			fileName = GenerateRandomString(15)
		}

		// Append path with file's name
		if strings.HasSuffix(outputPath, string(filepath.Separator)) {
			outputPath += GetLastPartOfPath(fileName)
		} else {
			outputPath += string(filepath.Separator) + GetLastPartOfPath(fileName)
		}
	} else {
		// Output is a file, create it's directory if needed
		os.MkdirAll(GetNthPartOfPath(outputPath, 1), 0750)
	}

	return outputPath
}
