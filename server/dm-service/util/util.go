package util

func Contains(list []string, target string) bool {
	for _, a := range list {
		if a == target {
			return true
		}
	}
	return false
}
