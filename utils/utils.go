package utils

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"

	"github.com/disintegration/imaging"
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

	originalsFolder := filepath.Join(threadDirPath, "o")
	err = os.Mkdir(originalsFolder, 0755)
	if err != nil {
		return err
	}
	thumbnailsFolder := filepath.Join(threadDirPath, "t")
	err = os.Mkdir(thumbnailsFolder, 0755)
	if err != nil {
		return err
	}

	return nil
}

func ValidatePost(title, text string) bool {
	return (title == "" && text != "") ||
		(title != "" && text == "") ||
		(title != "" && text != "")
}

func MakeImageThumbnail(originalPath, ext string, threadID, fileID int) (*image.NRGBA, error) {
	img, err := imaging.Open(originalPath, imaging.AutoOrientation(true))
	if err != nil {
		return nil, err
	}
	dstImage := imaging.Resize(img, 0, 150, imaging.NearestNeighbor)

	return dstImage, nil
}

func SaveImageThumbnail(img *image.NRGBA, threadID, fileID int, ext string) error {
	thumbnailPath := filepath.Join(
		UPLOADS_DIR_PATH,
		strconv.Itoa(threadID),
		"t",
		strconv.Itoa(fileID)+"."+ext,
	)

	f, err := os.Create(thumbnailPath)
	if err != nil {
		return err
	}

	switch ext {
	case "png":
		err = png.Encode(f, img)
	case "jpeg":
		err = jpeg.Encode(f, img, nil)
	}

	if err != nil {
		return err
	}

	return nil
}
