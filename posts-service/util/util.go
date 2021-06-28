package util

func Contains(s []*string, str string) bool {
	for _, v := range s {
		if *v == str {
			return true
		}
	}

	return false
}

func Search(s []*string, str string) int {
	for i, v := range s {
		if *v == str {
			return i
		}
	}

	return -1
}

func Compare(s1, s2 []*string) (result []*string) {
	for _, v := range s1 {
		if Contains(s2, *v) {
			result = append(result, v)
		}
	}
	return
}
