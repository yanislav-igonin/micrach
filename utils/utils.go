package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
)

var UPLOADS_DIR_PATH = "uploads"

// Check dir existence.
func CheckIfFolderExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Creates folder for uploads.
func CreateUploadsFolder() error {
	isExists := CheckIfFolderExists(UPLOADS_DIR_PATH)
	if isExists {
		return nil
	}

	err := os.Mkdir(UPLOADS_DIR_PATH, 0755)
	if err != nil {
		return err
	}

	return nil
}

// Creates folder for thread.
func CreateThreadFolder(postID int) error {

	threadDirPath := filepath.Join(UPLOADS_DIR_PATH, strconv.Itoa(postID))
	isExists := CheckIfFolderExists(threadDirPath)
	if isExists {
		return errors.New("folder already exists")
	}

	err := os.Mkdir(threadDirPath, 0755)
	if err != nil {
		return err
	}

	return nil
}

func ValidatePost(title, text string) bool {
	return (title == "" && text != "") || (title != "" && text == "")
}
