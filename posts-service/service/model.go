package service

type Post struct {
	ID			string 		`json:"id"` // username + date
	Description string 		`json:"description"`
	Data		string 		`json:"data"` // later filePath
	LikedBy		[]string 	`json:"likedBy"` // usernames
	Comments	[]string	`json:"comments"` // commentIds
}
