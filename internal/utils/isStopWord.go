package utils

import (
	"os"
	"path"
	"strings"
)

func IsStopWord(word string) (bool, error) {
	dirPath := "internal/data/stopwords"
	dir, err := os.ReadDir(dirPath)
	var isStopword bool = false

	if err != nil {
		return false, err
	}
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}

		data, err := os.ReadFile(path.Join(dirPath, entry.Name()))
		if err != nil {
			return false, err
		}
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			ass := strings.ToLower(strings.TrimSpace(line))
			if ass == strings.ToLower(word) {
				isStopword = true
				break
			}
		}
	}
	return isStopword, nil
}
