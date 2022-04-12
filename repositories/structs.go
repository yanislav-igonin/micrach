package repositories

import "time"

// DB Structs
// DB Structs
// DB Structs

type Post struct {
	ID        int       `json:"id"`
	IsParent  bool      `json:"-"`
	ParentID  *int      `json:"parentId"`
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
	return *p.ParentID
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
type Inputs struct {
	Title string
	Text  string
	Files string
}

// post-form.html
type HtmlFormData struct {
	FirstPostID     int
	CaptchaID       string
	IsCaptchaActive bool
	Errors          Inputs
	Inputs
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

const InvalidCaptchaErrorMessage = "Invalid captcha"
const InvalidTextOrFilesErrorMessage = "Text or files should not be empty"
const InvalidTitleLengthErrorMessage = "Title should not exceed 100 chars"
const InvalidTextLengthErrorMessage = "Text should not exceed 1000 chars"
const InvalidFilesLengthErrorMessage = "Maximum 4 files can be uploaded"
const InvalidFileSizeErrorMessage = "File size exceeded (3MB PER FILE)"
const InvalidFileExtErrorMessage = "Available file ext: PNG, JPG"
const ThreadIsArchivedErrorMessage = "Thread is archived"
