package repositories

import "time"

// DB Structs
// DB Structs
// DB Structs

type Post struct {
	ID        int       `json:"id"`
	IsParent  bool      `json:"-"`
	ParentID  int       `json:"parentId"`
	IsDeleted bool      `json:"-"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	IsSage    bool      `json:"isSage"`
	Files     []File    `json:"files"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"-"`
}

type File struct {
	ID        int       `json:"-"`
	PostID    int       `json:"-"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Ext       string    `json:"-"`
	Size      int       `json:"size"`
}

// HTML Templates Structs
// HTML Templates Structs
// HTML Templates Structs

// post-form.html
type HtmlFormData struct {
	FirstPostID int
	CaptchaID   string
}

// thread.html
type GetThreadHtmlData struct {
	Thread   []Post
	FormData HtmlFormData
}

// index.html
type GetThreadsHtmlData struct {
	Threads    []Post
	PagesCount int
	Page       int
	FormData   HtmlFormData
}

// 400.html
type BadRequestHtmlData struct {
	Message string
}

var InvalidTitleOrTextErrorMessage = "TITLE OR TEXT SHOULD NOT BE EMPTY"
var InvalidCaptchaErrorMessage = "INVALID CAPTCHA"
var InvalidFileSizeMessage = "FILE SIZE EXCIDED (3MB PER FILE)"
