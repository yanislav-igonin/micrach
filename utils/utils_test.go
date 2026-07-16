package utils

import (
	"mime/multipart"
	"net/textproto"
	"strings"
	"testing"

	"micrach/repositories"
)

func fileHeader(contentType string, size int64) *multipart.FileHeader {
	header := make(textproto.MIMEHeader)
	header.Set("Content-Type", contentType)
	return &multipart.FileHeader{Filename: "upload", Header: header, Size: size}
}

func TestValidatePost(t *testing.T) {
	validPNG := fileHeader("image/png", 1)
	validJPEG := fileHeader("image/jpeg", FILE_SIZE_IN_BYTES)

	tests := []struct {
		name  string
		title string
		text  string
		files []*multipart.FileHeader
		want  string
	}{
		{name: "text only", text: "hello"},
		{name: "file only", files: []*multipart.FileHeader{validPNG}},
		{name: "empty", want: repositories.InvalidTextOrFilesErrorMessage},
		{name: "title over limit", title: strings.Repeat("я", 101), text: "ok", want: repositories.InvalidTitleLengthErrorMessage},
		{name: "text over limit", text: strings.Repeat("я", 1001), want: repositories.InvalidTextLengthErrorMessage},
		{name: "four files", files: []*multipart.FileHeader{validPNG, validPNG, validJPEG, validJPEG}},
		{name: "five files", files: []*multipart.FileHeader{validPNG, validPNG, validPNG, validPNG, validPNG}, want: repositories.InvalidFilesLengthErrorMessage},
		{name: "unsupported content type", files: []*multipart.FileHeader{fileHeader("image/gif", 1)}, want: repositories.InvalidFileExtErrorMessage},
		{name: "size at limit", files: []*multipart.FileHeader{validJPEG}},
		{name: "size over limit", files: []*multipart.FileHeader{fileHeader("image/jpeg", FILE_SIZE_IN_BYTES+1)}, want: repositories.InvalidFileSizeErrorMessage},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidatePost(tt.title, tt.text, tt.files); got != tt.want {
				t.Fatalf("ValidatePost() = %q, want %q", got, tt.want)
			}
		})
	}
}
