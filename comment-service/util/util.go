package util

import "comment-service/graph/model"

func ConvertedIntoMapComments(currentComments map[string][]*model.Comment) []*model.MapComments {
	mapComments := make([]*model.MapComments, 0)

	for key, value := range currentComments {
		mapComment := &model.MapComments{
			Key:   key,
			Value: value,
		}

		mapComments = append(mapComments, mapComment)
	}

	return mapComments
}
