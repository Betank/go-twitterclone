package e2etests

type tweet struct {
	ID   string `json:"id"`
	User user   `json:"user"`
	Text string `json:"text"`
}

type user struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type stats struct {
	Follow   int `json:"follow"`
	Follower int `json:"follower"`
	Tweets   int `json:"tweets"`
}
