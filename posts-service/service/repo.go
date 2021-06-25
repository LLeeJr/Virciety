package service

import (
	"database/sql"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/lib/pq"
)

type Repository interface {
	CreatePost(post Post) error
	GetPosts() ([]Post, error)
	RemovePost(id string) (bool, error)
	EditPost(id, newDescription string) (bool, error)
	LikedPost(id, username string) (bool, error)
}

type repo struct {
	DB *sql.DB
	Posts []Post
	logger log.Logger
}

func NewRepo(db *sql.DB, logger log.Logger) (Repository, error) {
	return &repo{
		DB: db,
		Posts:  make([]Post, 0),
		logger: logger,
	}, nil
}

func (repo *repo) CreatePost(post Post) error {
	sqlQuery := `INSERT INTO posts (id, data, description, liked, comments) VALUES ($1, $2, $3, $4, $5)`

	_, err := repo.DB.Exec(sqlQuery, post.ID, post.Data, post.Description, pq.Array(post.LikedBy), pq.Array(post.Comments))
	if err != nil {
		level.Error(repo.logger).Log("CreatePostRepo ", err)
	}

	repo.logger.Log("Post Created: ", post.ID)

	return nil
}

func (repo *repo) GetPosts() ([]Post, error) {
	sqlQuery := `SELECT * FROM posts`

	posts := make([]Post, 0)

	rows, err := repo.DB.Query(sqlQuery)
	if err != nil {
		// handle this error better
		level.Error(repo.logger).Log("GetPostsRepoQuery ", err)
	}

	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Data, &post.Description, &post.ID, pq.Array(&post.LikedBy), pq.Array(&post.Comments))
		if err != nil {
			// handle this error better
			level.Error(repo.logger).Log("GetPostsRepoInForLoop ", err)
		}
		posts = append(posts, post)
	}
	// get any error encountered during this iteration
	err = rows.Err()
	if err != nil {
		level.Error(repo.logger).Log("GetPostsRepo ", err)
	}

	repo.logger.Log("Got all posts successfully")

	return posts, nil
}

func (repo *repo) RemovePost(id string) (bool, error) {
	sqlQuery := `DELETE FROM posts WHERE id = $1`

	res, err := repo.DB.Exec(sqlQuery, id)
	if err != nil {
		level.Error(repo.logger).Log("RemovePost ", err)
		return false, err
	}

	// validate how many rows been affected by upper query
	count, err := res.RowsAffected()
	if count != 1 {
		errMsg := "zero or more than one row has been removed"
		level.Error(repo.logger).Log("Remove post ", errMsg)
		return false, errors.New(errMsg)
	}

	repo.logger.Log("Deleted Post with id: ", id)

	return true, nil
}

func (repo *repo) EditPost(id, newDescription string) (bool, error) {
	sqlQuery := `UPDATE posts SET description = $1 WHERE id = $2`

	res, err := repo.DB.Exec(sqlQuery, newDescription, id)
	if err != nil {
		level.Error(repo.logger).Log("Edit post ", err)
		return false, err
	}

	// validate how many rows been affected by upper query
	count, err := res.RowsAffected()
	if count != 1 {
		errMsg := "zero or more than one row has been edited"
		level.Error(repo.logger).Log("Edit post ", errMsg)
		return false, errors.New(errMsg)
	}

	repo.logger.Log("Edited Post with id: ", id)

	return true, nil
}

func (repo *repo) LikedPost(id, username string) (bool, error) {
	// Get liked array out of database
	sqlQuery := `SELECT liked From posts WHERE id = $1`
	var liked []string

	row := repo.DB.QueryRow(sqlQuery, id)
	switch err := row.Scan(pq.Array(liked)); err {
	case sql.ErrNoRows:
		errMsg := "no rows were returned"
		level.Error(repo.logger).Log("Liked post select ", errMsg)
		return false, errors.New(errMsg)
	case nil:
		liked = append(liked, username)
	default:
		level.Error(repo.logger).Log("Liked post select ", err)
		return false, err
	}

	// update the liked array
	sqlQuery = `UPDATE posts SET liked = $1 WHERE id = $2`

	res, err := repo.DB.Exec(sqlQuery, pq.Array(liked), id)
	if err != nil {
		level.Error(repo.logger).Log("Liked post update ", err)
		return false, err
	}

	// validate how many rows been affected by upper query
	count, err := res.RowsAffected()
	if count != 1 {
		errMsg := "zero or more than one row has been edited"
		level.Error(repo.logger).Log("Liked post ", errMsg)
		return false, errors.New(errMsg)
	}

	return true, nil
}