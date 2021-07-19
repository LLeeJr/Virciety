package database

type PostEvent struct {
	EventTime   string   `json:"event_time"`
	EventType   string   `json:"event_type"`
	PostID      string   `json:"id"`
	Username    string   `json:"username"`
	Description string   `json:"description"`
	Data        string   `json:"data"`
	LikedBy     []string `json:"liked_by"`
	Comments    []string `json:"comments"`
}

type CommentEvent struct {
	EventTime   string   `json:"event_time"`
	EventType   string   `json:"event_type"`
	CommentID   string   `json:"id"`
	Username    string   `json:"username"`
	Description string   `json:"description"`
	LikedBy     []string `json:"liked_by"`
	PostID      string   `json:"post_id"`
}
