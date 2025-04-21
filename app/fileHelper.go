package main

import (
	"io"
	"os"
)

func GetFileContent(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return content, nil
}
func AddFile(fileName string, content string) bool {
	file, err := os.Create(fileName)
	if err != nil {
		return false
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return false
	}
	return true
}
