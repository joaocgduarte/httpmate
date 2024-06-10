package files

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func GetFilesFromDirectoryWithoutExtension(parentDirectory string) []string {
	jsonFiles := []string{}
	err := filepath.WalkDir(parentDirectory, func(path string, d os.DirEntry, err error) error {
		cobra.CheckErr(err)
		if parentDirectory == path {
			return nil
		}

		if !d.IsDir() {
			path = path[:len(path)-len(filepath.Ext(path))]
			path = strings.Replace(path, parentDirectory, "", -1)
			jsonFiles = append(jsonFiles, path)
		}
		return nil
	})

	cobra.CheckErr(err)
	return jsonFiles
}

func GetSubDirectories(parentDirectory string) []string {
	result := make([]string, 0)
	err := filepath.WalkDir(parentDirectory, func(path string, d os.DirEntry, err error) error {
		cobra.CheckErr(err)
		if parentDirectory == path {
			return nil
		}

		if d.IsDir() {
			path = path[:len(path)-len(filepath.Ext(path))]
			path = strings.Replace(path, parentDirectory, "", -1)
			result = append(result, path)
		}
		return nil
	})
	cobra.CheckErr(err)
	return result
}

func CreateDirectory(directory string) {
	err := os.MkdirAll(directory, 0755)
	cobra.CheckErr(err)
}

func WriteStructToJSONFile(data interface{}, filepath string) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	cobra.CheckErr(err)

	file, err := os.Create(filepath)
	cobra.CheckErr(err)
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	cobra.CheckErr(err)
}
