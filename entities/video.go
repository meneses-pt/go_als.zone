package entities

// Video is the representation entity for videos
type Video struct {
	Title      string    `json:"title"`
	RedditLink string    `json:"reddit_link"`
	Mirrors    *[]Mirror `json:"mirrors"`
}
