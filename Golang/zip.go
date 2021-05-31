package goutils

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

// ZipFilesFromPaths packs every file contained in the slice into a single zip file
func ZipFilesFromPaths(filePaths []string) (*[]byte, error) {

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	zipWriter := zip.NewWriter(buf)

	// Add some files to the archive.
	for _, file := range filePaths {
		// Write the file or all the dir's content
		err := writeFileOrDirectoryToZip(file, zipWriter)
		if err != nil {
			return nil, err
		}
	}

	// Closing and checking for errors
	err := zipWriter.Close()
	if err != nil {
		return nil, err
	}

	// Return the bytes
	byts := buf.Bytes()
	buf.Reset()

	return &byts, nil
}

func writeFileOrDirectoryToZip(filePath string, w *zip.Writer) error {
	// Open the file
	openedFile, err := os.Open(filePath)
	if err != nil {
		return errors.New("Couldn't read File: " + filePath + " | " + err.Error())
	}

	// Given "file" is a directory
	info, err := openedFile.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {

		// Read every file
		files, err := ioutil.ReadDir(filePath)
		if err != nil {
			return errors.New("Couldn't read Dir: " + filePath + " | " + err.Error())
		}
		for _, f := range files {
			innerFilePath := filePath + string(filepath.Separator) + f.Name()

			if f.IsDir() {
				// Inner File is also a Dir
				err = writeFileOrDirectoryToZip(innerFilePath, w)
				if err != nil {
					return err
				}
			} else {
				err = writeFileToZip(innerFilePath, filePath, w)
				if err != nil {
					return err
				}
			}
		}
	} else {
		// Given File is really just a file
		writeFileToZip(filePath, filePath, w)
	}

	return nil
}

func writeFileToZip(innerFilePath string, filePath string, w *zip.Writer) error {
	// Remove relative part components
	innerFilePath = ReplaceFilepathSeparator(innerFilePath, string(filepath.Separator))
	innerFilePath = strings.ReplaceAll(innerFilePath, ".."+string(filepath.Separator), "")
	// Write File to Zip
	zipFile, err := w.Create(innerFilePath)
	if err != nil {
		return errors.New("Couldnt create File in zip: " + err.Error())
	}
	zipFileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return errors.New("Couldnt read File: " + filePath + " | " + err.Error())
	}
	zipFile.Write(zipFileContent)
	return err
}

// CreateTarFromDirectory creates a tar file from the given Dir
func CreateTarFromDirectory(src string, buf io.Writer) error {
	maxErrors := 10

	tw := tar.NewWriter(buf)
	buff := make([]byte, 1024*1024)
	// baseDir := getBaseDir(src)

	errChan := make(chan error, maxErrors)

	// walk through every file in the folder
	go func() {
		filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
			if len(file) < len(src)+1 {
				return nil
			}

			// Follow link
			var link string
			if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
				if link, err = os.Readlink(file); err != nil {
					errChan <- err
					return nil
				}
			}

			// Generate tar header
			header, err := tar.FileInfoHeader(fi, link)
			if err != nil {
				errChan <- err
				return nil
			}

			// Set filename
			header.Name = filepath.Join(src, strings.TrimPrefix(file, src))
			//header.Name = filepath.ToSlash(file)

			// write header
			if err := tw.WriteHeader(header); err != nil {
				errChan <- err
				return nil
			}

			// Nothing more to do for non-regular
			if !fi.Mode().IsRegular() {
				return nil
			}

			// can only write file-
			// contents to archives
			if !fi.IsDir() {
				data, err := os.Open(file)
				if err != nil {
					errChan <- err
					return nil
				}

				if _, err := io.CopyBuffer(tw, data, buff); err != nil {
					errChan <- err
					return nil
				}

				data.Close()
			}

			return nil
		})

		close(errChan)
	}()

	errCounter := 0
	for err := range errChan {
		if errCounter >= maxErrors {
			return errors.New("Too many errors")
		}

		fmt.Println(err)
		errCounter++
	}

	// produce tar
	if err := tw.Close(); err != nil {
		return err
	}

	return nil
}
