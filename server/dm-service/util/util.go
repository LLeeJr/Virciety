package util

import "fmt"

func Contains(list []string, target string) bool {
	for _, a := range list {
		fmt.Println("USER", a)
		if a == target {
			return true
		}
	}
	return false
}
