package goutils

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/JojiiOfficial/shred"
)

// DoesFileExist checks íf a given file exists or not
func DoesFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}

// ReadFileToString returns a string with the file's content. Returns "" on error
func ReadFileToString(filePath string) string {
	var fileBytes []byte
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return ""
	}
	return string(fileBytes)
}

// ShredderFile shreddres a given file
func ShredderFile(localFile string, size int64) {
	shredder := shred.Shredder{}

	var shredConfig *shred.ShredderConf
	if size < 0 {
		s, err := os.Stat(localFile)
		if err != nil {
			fmt.Println("File to shredder not found")
			return
		}
		size = s.Size()
	}

	if size >= 1000000000 {
		// Size >= 1GB
		shredConfig = shred.NewShredderConf(&shredder, shred.WriteZeros, 1, true)
	} else if size >= 1000000000 {
		// Size >= 1GB
		shredConfig = shred.NewShredderConf(&shredder, shred.WriteZeros|shred.WriteRandSecure, 2, true)
	} else if size >= 5000 {
		// Size > 5kb
		shredConfig = shred.NewShredderConf(&shredder, shred.WriteZeros|shred.WriteRandSecure, 3, true)
	} else {
		// Size < 5kb
		shredConfig = shred.NewShredderConf(&shredder, shred.WriteZeros|shred.WriteRandSecure, 6, true)
	}

	// Shredder & Delete local file
	err := shredConfig.ShredFile(localFile)
	if err != nil {
		fmt.Println(err)
		// Delete file if shredder didn't
		err = os.Remove(localFile)
		if err != nil {
			fmt.Println(err)
		}
	}
}
