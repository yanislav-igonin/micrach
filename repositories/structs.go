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

func (p Post) GetThreadID() int {
	if p.IsParent {
		return p.ID
	}
	return p.ParentID
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
	FirstPostID     int
	CaptchaID       string
	IsCaptchaActive bool
}

// index.html
type HtmlPaginationData struct {
	PagesCount int
	Page       int
}

// thread.html
type GetThreadHtmlData struct {
	Thread   []Post
	FormData HtmlFormData
}

// index.html
type GetThreadsHtmlData struct {
	Threads    []Post `json:"threads"`
	Pagination HtmlPaginationData
	FormData   HtmlFormData
}

// 400.html
type BadRequestHtmlData struct {
	Message string
}

const InvalidCaptchaErrorMessage = "INVALID CAPTCHA"
const InvalidTextOrFilesErrorMessage = "TEXT OR FILES SHOULD NOT BE EMPTY"
const InvalidTitleLengthErrorMessage = "TITLE SHOULD NOT EXCEED 100 CHARS"
const InvalidTextLengthErrorMessage = "TEXT SHOULD NOT EXCEED 1000 CHARS"
const InvalidFilesLengthErrorMessage = "MAXIMUM 4 FILES CAN BE UPLOADED"
const InvalidFileSizeErrorMessage = "FILE SIZE EXCIDED (3MB PER FILE)"
const InvalidFileExtErrorMessage = "AVALIABLE FILE EXT: PNG, JPG"
const ThreadIsArchivedErrorMessage = "THREAD IS ARCHIVED"
