package command

import (
	"os"
	"path"

	"golang.org/x/mod/modfile"
)

func Module(dir string) (string, error) {
	if dir == "" {
		dir = "."
	}
	filename := path.Join(dir, "go.mod")
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	file, err := modfile.Parse(filename, data, nil)
	if err != nil {
		return "", err
	}
	return file.Module.Mod.Path, nil
}
