package util

// basic slice operations

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func Search(s []string, str string) int {
	for i, v := range s {
		if v == str {
			return i
		}
	}

	return -1
}
