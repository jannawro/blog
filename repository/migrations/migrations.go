package migrations

import (
	"embed"
	"io/fs"
)

//go:embed *.sql
var migrationsFiles embed.FS

func Files() fs.FS {
	return migrationsFiles
}

func DebugListFiles() []string {
	var files []string
	err := fs.WalkDir(migrationsFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}
