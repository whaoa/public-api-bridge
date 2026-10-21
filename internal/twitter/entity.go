package twitter

import "time"

const (
	MediaTypePhoto = "photo"
	MediaTypeVideo = "video"
)

type User struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ScreenName  string    `json:"screen_name"`
	Avatar      string    `json:"avatar"`
	Banner      string    `json:"banner"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	JoinedAt    time.Time `json:"joined_at"`
	Verified    bool      `json:"verified"`
	Protected   bool      `json:"protected"`
	Following   int64     `json:"following"`
	Followers   int64     `json:"followers"`
	Url         string    `json:"url"`
}

type Tweet struct {
	ID                string        `json:"id"`
	Author            *BasicUser    `json:"author"`
	Content           string        `json:"content"`
	Attachments       []*Attachment `json:"attachments"`
	MosaicPhoto       string        `json:"mosaic_photo"`
	Replies           int64         `json:"replies"`
	Retweets          int64         `json:"retweets"`
	Likes             int64         `json:"likes"`
	Views             int64         `json:"views"`
	Quote             *BasicTweet   `json:"quote"`
	ReplyingTo        *BasicTweet   `json:"replying_to"`
	CreatedAt         time.Time     `json:"created_at"`
	PossiblySensitive bool          `json:"possibly_sensitive"`
	Url               string        `json:"url"`
}

type BasicUser struct {
	ScreenName string `json:"screen_name"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
	Url        string `json:"url"`
}

type Attachment struct {
	Kind      string `json:"type"`
	Url       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
}

type BasicTweet struct {
	ScreenName string `json:"screen_name"`
	Tweet      string `json:"tweet"`
	Url        string `json:"url"`
}

type TweetSlice struct {
	Items      []*Tweet `json:"items"`
	PrevCursor string   `json:"prev_cursor"`
	NextCursor string   `json:"next_cursor"`
}
