package files

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func GetFilesInFolder(folder string) []string {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(filepath.Join(currentPath + "/migrations"))
	if err != nil {
		log.Fatal(err)
	}

	var filesNames []string

	for _, file := range files {
		filesNames = append(filesNames, file.Name())
	}

	return filesNames
}
