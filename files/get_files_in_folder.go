package files

import (
	"io/ioutil"
	"log"
	"os"
	"path"
)

// Reads folder and returns full file paths slice
func GetFullFilePathsInFolder(folder string) []string {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fullFolderPath := path.Join(currentPath, folder)
	files, err := ioutil.ReadDir(fullFolderPath)
	if err != nil {
		log.Fatal(err)
	}

	var paths []string

	for _, file := range files {
		paths = append(paths, path.Join(fullFolderPath, file.Name()))
	}

	return paths
}

// Reads file contents by full path and returns string
func ReadFileText(fullFilePath string) string {
	file, err := ioutil.ReadFile(fullFilePath)
	if err != nil {
		log.Fatal(err)
	}

	return string(file)
}
