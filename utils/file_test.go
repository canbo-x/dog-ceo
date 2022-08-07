package utils

import (
	"os"
	"testing"
)

func TestSaveToDisk(t *testing.T) {
	tests := map[string]struct {
		FileName         string
		FilePath         string
		ExpectedFullPath string
		Data             []byte
		Valid            bool
	}{
		"valid file name": {
			FileName:         "test.txt",
			FilePath:         "/tmp",
			ExpectedFullPath: "/tmp/test.txt",
			Data:             []byte("test"),
			Valid:            true,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			fullPath, err := SaveToDisk(test.Data, test.FileName, test.FilePath)
			if (err == nil) != test.Valid {
				t.Fatalf("want err == nil => %t; got err %v", test.Valid, err)
			}
			if fullPath != test.ExpectedFullPath {
				t.Fatalf("full path is not correct got: %v expected: %v", fullPath, test.ExpectedFullPath)
			}
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				t.Fatalf("file does not exist")
			}
			if err := os.Remove(fullPath); err != nil {
				t.Errorf("Error is not nil %v", err)
			}
		})
	}

}

func TestExists(t *testing.T) {
	tests := map[string]struct {
		FilePath    string
		shouldExist bool
	}{
		"valid file path": {
			FilePath:    "/tmp/test_this_file_exists.txt",
			shouldExist: true,
		},
		"invalid file path": {
			FilePath:    "/not/a/real/path",
			shouldExist: false,
		},
	}

	fullPath, err := SaveToDisk([]byte("test"), "test_this_file_exists.txt", "/tmp")
	if err != nil {
		t.Errorf("Error is not nil %v", err)
	}

	defer func() {
		if err := os.Remove(fullPath); err != nil {
			t.Errorf("Error is not nil %v", err)
		}
	}()

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			doesExist, err := exists(test.FilePath)
			if err != nil {
				t.Errorf("Error is not nil %v", err)
			}
			if doesExist != test.shouldExist {
				t.Errorf("expected %t got %t", test.shouldExist, doesExist)
			}
		})
	}
}
