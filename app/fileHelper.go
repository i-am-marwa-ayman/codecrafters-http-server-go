package main

import (
	"io"
	"os"
)

func GetFileContent(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
func AddFile(fileName string, content string) bool {
	file, err := os.Create(fileName)
	defer file.Close()

	if err != nil {
		return false
	}
	_, err = file.WriteString(content)
	if err != nil {
		return false
	}
	return true
}
