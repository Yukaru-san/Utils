package goutils

// Sort2DByteSliceByLength sorts the given slice by length of inner bytes
func Sort2DByteSliceByLength(slice *[][]byte) {
	for i := len(*slice); i > 0; i-- {
		for j := 1; j < i; j++ {
			if len((*slice)[j-1]) < len((*slice)[j]) {
				tmp := (*slice)[j]
				(*slice)[j] = (*slice)[j-1]
				(*slice)[j-1] = tmp
			}
		}
	}
}
