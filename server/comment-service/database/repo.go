package database

import (
	"comment-service/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateComment(event CommentEvent) (*model.Comment, error)
	GetComments() (map[string][]*model.Comment, error)
	GetCommentsByPostId(postId string) ([]*model.Comment, error)
	GetCommentById(commentId string) (*model.Comment, int, string, error)
	RemoveComment(event CommentEvent, index int) (string, error)
	EditComment(event CommentEvent) (string, error)
	LikeComment(event CommentEvent) (string, error)
	UnlikeComment(event CommentEvent) (string, error)
}

type repo struct {
	commentCollection *mongo.Collection
}

func NewRepo() (Repository, error) {
	client, err := dbConnect()
	if err != nil {
		return nil, err
	}

	db := client.Database("comment-service")

	return &repo{
		commentCollection: db.Collection("comment-events"),
	}, nil
}

func (repo *repo) InsertCommentEvent(commentEvent CommentEvent) (string, error) {
	inserted, err := repo.commentCollection.InsertOne(ctx, commentEvent)
	if err != nil {
		return "", err
	}

	return inserted.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (repo *repo) CreateComment(event CommentEvent) (*model.Comment, error) {
	insertedID, err := repo.InsertCommentEvent(event)
	if err != nil {
		return nil, err
	}

	comment := &model.Comment{
		ID:        insertedID,
		PostID:    event.PostID,
		Comment:   event.Comment,
		CreatedBy: event.CreatedBy,
	}

	return comment, nil
}

func (repo *repo) GetComments() (map[string][]*model.Comment, error) {
	/*currentComments := make([]*model.Comment, 0)

		// first get all rows with event_type = "CreateComment" and latestEventId
		sqlQuery := `select "commentID", description, liked, "postID", id from "comment-events" where id > $1 and "eventType" = $2 `

		rows, err := repo.DB.Query(sqlQuery, repo.currentEventId, "CreateComment")
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		oldId := repo.currentEventId
		id := repo.currentEventId

		for rows.Next() {
			var comment model.Comment

			err = rows.Scan(&comment.ID, &comment.Description, pq.Array(&comment.LikedBy), &comment.PostID, &id)
			if err != nil {
				repo.currentEventId = oldId
				return nil, err
			}
			currentComments = append(currentComments, &comment)
		}

		// list is recent
		if id == repo.currentEventId {
			log.Printf("Comment list is up to date!")
			return repo.currentComments, nil
		}

		repo.currentEventId = id

		for _, comment := range currentComments {
			sqlQuery = `select liked, description from "comment-events" where id = (select max(id) from "comment-events" where "commentID" = $1
	                                                                                                   and ("eventType" = $2 or "eventType" = $3 or "eventType" = $4))`

			row := repo.DB.QueryRow(sqlQuery, comment.ID, "EditComment", "LikeComment", "UnlikeComment")

			switch err = row.Scan(pq.Array(&comment.LikedBy), &comment.Description); err {
			case sql.ErrNoRows:
				// nothing happens because it is not really an error
				// since a post doesn't have to be edited
			case nil:
				log.Printf("Edited data added to " + comment.ID)
			default:
				repo.currentEventId = oldId
				return nil, err
			}

			// add to currentComments
			comments, ok := repo.currentComments[comment.PostID]

			if ok {
				repo.currentComments[comment.PostID] = append(comments, comment)
			} else {
				repo.currentComments[comment.PostID] = []*model.Comment{comment}
			}
		}*/

	return nil, nil
}

func (repo *repo) GetCommentsByPostId(postId string) ([]*model.Comment, error) {
	/*comments, ok := repo.currentComments[postId]
	if !ok {
		errMsg := "no comments for post with id " + postId + " found"
		return nil, errors.New(errMsg)
	}

	return comments, nil*/
	return nil, nil
}

func (repo *repo) GetCommentById(commentId string) (*model.Comment, int, string, error) {
	/*// process data to get postId
	info := strings.Split(commentId, "__")
	index := -1

	username := info[1]
	postId := info[2] + "__" + info[3]

	// get comment data out of current saved comments
	comments, err := repo.GetCommentsByPostId(postId)
	if err != nil {
		return nil, index, username, err
	}

	// search comment in post comments list
	var comment *model.Comment
	for i, v := range comments {
		if v.ID == commentId {
			comment = v
			index = i
			break
		}
	}

	if comment == nil {
		errMsg := "no comment with id " + commentId + " for post with id " + postId + " found"
		return nil, index, username, errors.New(errMsg)
	}

	return comment, index, username, nil*/
	return nil, 0, "", nil
}

func (repo *repo) RemoveComment(event CommentEvent, index int) (string, error) {
	/*// remove from currentComments
	comments, err := repo.GetCommentsByPostId(event.PostID)
	if err != nil {
		return "failed", err
	}

	comments = append(comments[:index], comments[index+1:]...)

	if len(comments) == 0 {
		delete(repo.currentComments, event.PostID)
	} else {
		repo.currentComments[event.PostID] = comments
	}

	// delete all events relating to the id
	sqlQuery := `delete from "comment-events" where "commentID" = $1`

	_, err = repo.DB.Exec(sqlQuery, event.CommentID)
	if err != nil {
		return "failed", err
	}

	err = repo.InsertCommentEvent(event)
	if err != nil {
		return "failed", err
	}
	*/
	return "success", nil
}

func (repo *repo) EditComment(event CommentEvent) (string, error) {
	_, err := repo.InsertCommentEvent(event)
	if err != nil {
		return "failed", err
	}

	return "success", nil
}

func (repo *repo) LikeComment(event CommentEvent) (string, error) {
	_, err := repo.InsertCommentEvent(event)
	if err != nil {
		return "failed", err
	}

	return "success", nil
}

func (repo *repo) UnlikeComment(event CommentEvent) (string, error) {
	_, err := repo.InsertCommentEvent(event)
	if err != nil {
		return "failed", err
	}

	return "success", nil
}
