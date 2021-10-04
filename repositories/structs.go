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

type HtmlFormData struct {
	FirstPostID int
	CaptchaID   string
}

type GetThreadHtmlData struct {
	Thread   []Post
	FormData HtmlFormData
}

type GetThreadsHtmlData struct {
	Threads    []Post
	PagesCount int
	Page       int
	FormData   HtmlFormData
}
