package goutils

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// FindContentBetween tries to find the first occurence of some text between start and end sequence. Returns "", false when nothing was found
func FindContentBetween(baseContent string, startSequence string, endSequence string) (string, bool) {
	startIndex := strings.Index(baseContent, startSequence)
	if startIndex == -1 {
		return "", false
	}
	startIndex += len(startSequence)

	endIndex := strings.Index(baseContent[startIndex:], endSequence) + startIndex
	if endIndex == -1 {
		return "", false
	}

	baseContent = baseContent[startIndex:endIndex]
	return baseContent, true
}

// SplitStringByLine splits a string by lines and returns them as a slice
func SplitStringByLine(base string) []string {
	return strings.Split(strings.ReplaceAll(base, "\r\n", "\n"), "\n")
}

// GenerateRandomString creates a random string of length n
func GenerateRandomString(n int) string {

	// Seed the time
	rand.Seed(time.Now().UnixNano())

	// Available Chars
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	// Generate
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// StrToUIntslice converts an string slice to an integer slice
func StrToUIntslice(sl []string) ([]uint64, error) {
	iSl := make([]uint64, len(sl))
	for i := range sl {
		var err error
		iSl[i], err = strconv.ParseUint(sl[i], 10, 32)
		if err != nil {
			return nil, err
		}
	}

	return iSl, nil
}

// ArgsToDate takes an arg like 2d and turns int into something readable
// returns an error if the input was faulty
func ArgsToDate(args []string) (time.Time, error) {

	// Get the current date and ignore seconds
	endDate := time.Now()
	endDate = endDate.Add(time.Second * time.Duration(-1*endDate.Second()))
	fmt.Println(endDate)

	// Direct date calling
	if strings.Contains(args[0], ".") || strings.Contains(args[0], ":") {

		// Day, Month, Year, Hour, Minute
		timeInformation := make([]int64, 5)

		// Initially fill with current data
		timeInformation[0] = int64(endDate.Day())
		timeInformation[1] = int64(endDate.Month())
		timeInformation[2] = int64(endDate.Year())
		timeInformation[3] = int64(endDate.Hour())
		timeInformation[4] = int64(endDate.Minute())

		for _, arg := range args {
			// Direct date?
			if strings.Contains(arg, ".") {

				// Directly called Day, Month and/or Year
				dates := strings.Split(arg, ".")

				for i, d := range dates {

					// Extract the date the user wanted
					dateWanted, err := strconv.ParseInt(d, 10, 64)

					// Input couldn't be parsed
					if err != nil {
						return endDate, err
					}

					// Add to the data slice
					timeInformation[i] = dateWanted
				}
			} else if strings.Contains(arg, ":") {
				// Directly called end Time

				// expecting 2 results in slice
				dates := strings.Split(arg, ":")
				if len(dates) != 2 {
					return endDate, errors.New("A specific time has to contain hh:mm")
				}

				// Extract the date the user wanted
				hour, err := strconv.ParseInt(dates[0], 10, 64)
				minute, err := strconv.ParseInt(dates[1], 10, 64)

				// Input couldn't be parsed
				if err != nil {
					return endDate, err
				}

				// Add informations
				timeInformation[3] = hour
				timeInformation[4] = minute
			}
		} // Day, Month, Year, Hour, Minute

		// Parse all the gathered informations
		endDate = time.Date(int(timeInformation[2]), time.Month(timeInformation[1]), int(timeInformation[0]), int(timeInformation[3]), int(timeInformation[4]), 0, 0, time.UTC)

	} else {
		// Intuitive date calling
		for _, arg := range args {

			// Extract the days the user wanted
			daysWanted, err := strconv.ParseInt(arg[:len(arg)-1], 10, 64)

			// Input couldn't be parsed
			if err != nil {
				return endDate, err
			}

			// Adjust Time
			if strings.HasSuffix(arg, "d") {
				endDate = endDate.Add(time.Hour * 24 * time.Duration(daysWanted))
			} else if strings.HasSuffix(arg, "h") {
				endDate = endDate.Add(time.Hour * time.Duration(daysWanted))
			} else if strings.HasSuffix(arg, "m") {
				endDate = endDate.Add(time.Minute * time.Duration(daysWanted))
			}
		}
	}

	return endDate, nil
}

// ReplaceAndGetIndex is literally the Replace function from the strings lib but also returns index of replacements
func ReplaceAndGetIndex(s, old, new string, n int) (string, []int) {
	if old == new || n == 0 {
		return s, []int{-1} // avoid allocation
	}

	// Compute number of replacements.
	if m := strings.Count(s, old); m == 0 {
		return s, []int{-1} // avoid allocation
	} else if n < 0 || m < n {
		n = m
	}

	// Position counter
	c := make([]int, n)

	// Apply replacements to buffer.
	t := make([]byte, len(s)+n*(len(new)-len(old)))
	w := 0
	start := 0
	for i := 0; i < n; i++ {
		j := start
		if len(old) == 0 {
			if i > 0 {
				_, wid := utf8.DecodeRuneInString(s[start:])
				j += wid
			}
		} else {
			j += strings.Index(s[start:], old)
		}
		w += copy(t[w:], s[start:j])
		w += copy(t[w:], new)

		c[i] = j - i*len(old)
		start = j + len(old)
	}
	w += copy(t[w:], s[start:])
	return string(t[0:w]), c
}
