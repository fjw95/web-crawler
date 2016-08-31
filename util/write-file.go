package util

import (
	"os"
	"path/filepath"
)

func WriteFile(content_file []string, path_file string) {

	absPath, _ := filepath.Abs(path_file)

	f, err := os.Create(absPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for _, url := range content_file {
		f.WriteString(url + "\n")
	}

	f.Sync()
}
