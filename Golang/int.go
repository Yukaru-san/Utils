package goutils

// IsContainedInSlice checks if the given value is contained within the slice
func IsContainedInSlice(slice *[]int, value int) bool {
	for _, e := range *slice {
		if e == value {
			return true
		}
	}

	return false
}

// ToUniqueSlice takes and integer slice and removes dublicates
func ToUniqueSlice(input []int) []int {
	keys := make(map[int]bool)
	uniqueList := []int{}
	for _, entry := range input {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueList = append(uniqueList, entry)
		}
	}

	return uniqueList
}
