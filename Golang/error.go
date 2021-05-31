package goutils

import "fmt"

// CheckErr checks for an error and prints if found
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// CheckErrAndPanic checks for an error and panics if found
func CheckErrAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}
