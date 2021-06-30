package util

import (
	"posts-service/graph/model"
)

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

// Useful stuff

func GetPostByID(currentPosts []*model.Post, id string) (int, *model.Post) {
	for i, post := range currentPosts {
		if post.ID == id {
			return i, post
		}
	}

	return -1, nil
}
