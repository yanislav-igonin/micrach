package build

import (
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
)

// Gets file paths from directory recursively.
func getFilePathsRecursively(dir string) ([]string, error) {
	var paths []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			subdir := dir + "/" + file.Name()
			subpaths, err := getFilePathsRecursively(subdir)
			if err != nil {
				return nil, err
			}
			paths = append(paths, subpaths...)
			continue
		}
		paths = append(paths, dir+"/"+file.Name())
	}
	return paths, nil
}

// Returns all css files paths in ../static/styles folder.
func getCssFilesPaths() []string {
	var paths []string
	files, err := ioutil.ReadDir("static/styles")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		paths = append(paths, "static/styles/"+file.Name())
	}
	return paths
}

// Creates map of css file names to their new names.
func createCssMap() map[string]string {
	var cssMap = make(map[string]string)
	origPaths := getCssFilesPaths()
	for _, origPath := range origPaths {
		newPath := "static/styles/" + randomString(10) + ".css"
		cssMap[origPath] = newPath
	}

	return cssMap
}

// Renames file by paths.
func renameFile(oldPath, newPath string) {
	// rename the file
	err := os.Rename(removeFirstChar(oldPath), removeFirstChar(newPath))
	if err != nil {
		panic(err)
	}
}

// Removes first character from string.
func removeFirstChar(s string) string {
	return s[1:]
}

// Generates a random string.
func randomString(length int) string {
	// create a slice of characters to use
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	// create a new slice to hold the random string
	b := make([]rune, length)

	// loop through the length of the string
	for i := range b {
		// get a random number
		b[i] = letters[rand.Intn(len(letters))]
	}

	// return the string
	return string(b)
}

// Renames css files to prevent caching old files on production.
func RenameCss() {
	cssMap := createCssMap()
	htmlTemplatePaths, err := getFilePathsRecursively("templates")
	if err != nil {
		panic(err)
	}

	for origPath, newPath := range cssMap {
		renameFile(origPath, newPath)
	}

	for _, htmlTemplatePath := range htmlTemplatePaths {
		htmlTemplate, err := ioutil.ReadFile(htmlTemplatePath)
		if err != nil {
			panic(err)
		}
		// walk through css map and replace all occurencies of css origName with newName
		for origPath, newPath := range cssMap {
			htmlTemplate = []byte(strings.Replace(string(htmlTemplate), origPath, newPath, -1))
		}
		// write html template to file
		err = ioutil.WriteFile(htmlTemplatePath, htmlTemplate, 0644)
		if err != nil {
			panic(err)
		}
	}
}
