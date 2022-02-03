package util

func Contains(list []string, target string) bool {
	for _, a := range list {
		if a == target {
			return true
		}
	}
	return false
}

func Remove(list []string, target string) []string {
	var index int
	for i, s := range list {
		if s == target {
			index = i
		}
	}
	return append(list[:index], list[index+1:]...)
}