package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// This is the default file permission in order to write a file.
const defaultFilePerm os.FileMode = 0666

// This is the default folder permission in order to create a directory
const defaultFolderPerm os.FileMode = 0777

// SaveToDisk saves the image to the disk with the given file path and file name and returns the saved file path and an error if any.
func SaveToDisk(image []byte, fileName, filePath string) (string, error) {
	// TODO : handle concurrency because os file operations are not thread safe.
	exist, err := exists(filePath)
	if err != nil {
		return "", fmt.Errorf("file path error : %v", err)
	}
	if !exist {
		if os.MkdirAll(filePath, defaultFolderPerm); err != nil {
			return "", fmt.Errorf("failed to create directory : %v", err)
		}
	}

	fullPath := filepath.Join(filePath, fileName)
	if err := os.WriteFile(fullPath, image, defaultFilePerm); err != nil {
		return "", fmt.Errorf("failed to write file : %v", err)
	}

	return fullPath, nil
}

// exists returns whether the given file or directory exists and an error if any.
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
