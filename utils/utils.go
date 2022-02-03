// TODO: move all functions to different packages
package utils

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	Repositories "micrach/repositories"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/disintegration/imaging"
)

type stringSlice []string

const UPLOADS_DIR_PATH = "uploads"
const FILE_SIZE_IN_BYTES = 3145728                               // 3MB
var PERMITTED_FILE_EXTS = stringSlice{"image/jpeg", "image/png"} // 3MB

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

func ValidatePost(title, text string, files []*multipart.FileHeader) string {
	if text == "" && len(files) == 0 {
		return Repositories.InvalidTextOrFilesErrorMessage
	}

	if len([]rune(title)) > 100 {
		return Repositories.InvalidTitleLengthErrorMessage
	}

	if len([]rune(text)) > 1000 {
		return Repositories.InvalidTextLengthErrorMessage
	}

	if len(files) > 4 {
		return Repositories.InvalidFilesLengthErrorMessage
	}

	isFilesExtsValid := CheckFilesExt(files)
	if !isFilesExtsValid {
		return Repositories.InvalidFileExtErrorMessage
	}

	isFilesSizesNotToBig := CheckFilesSize(files)
	if !isFilesSizesNotToBig {
		return Repositories.InvalidFileSizeErrorMessage
	}

	return ""
}

func CheckFilesSize(files []*multipart.FileHeader) bool {
	for _, file := range files {
		if file.Size > int64(FILE_SIZE_IN_BYTES) {
			return false
		}
	}

	return true
}

func CheckFilesExt(files []*multipart.FileHeader) bool {
	for _, file := range files {
		ext := file.Header.Get("Content-Type")
		if !PERMITTED_FILE_EXTS.includes(ext) {
			return false
		}
	}

	return true
}

func (ss stringSlice) includes(toCheck string) bool {
	for _, s := range ss {
		if toCheck == s {
			return true
		}
	}
	return false
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
