package model


type Dcard struct {
	ID        int              `json:"id"`
	Media     []*DcardImageUrl `json:"mediaMeta"`
	Gender    string           `json:"gender"`
	LikeCount int              `json:"likeCount"`
	Title     string           `json:"title"`
}

type DcardInfo struct {
	ID    int    `json:"id"`
	Image string `json:"image"`
	Link  string `json:"link"`
	Title string `json:"title"`
}

type DcardImageUrl struct {
	Url string `json:"url"`
}

