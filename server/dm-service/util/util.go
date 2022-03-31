package util

// Contains is a list helper-function for checking if a list contains a target-value or not
func Contains(list []string, target string) bool {
	for _, a := range list {
		if a == target {
			return true
		}
	}
	return false
}

// Remove is a list helper-function for removing target-value from a list (provided it is contained)
func Remove(list []string, target string) []string {
	var index int
	for i, s := range list {
		if s == target {
			index = i
		}
	}
	return append(list[:index], list[index+1:]...)
}