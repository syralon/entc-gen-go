package text

import (
	"os"
	"path"
	"strings"

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

func ProtoModule(module string) string {
	temp := strings.Split(module, "/")
	if len(temp) < 3 {
		return path.Join("proto", path.Base(module))
	}
	return path.Join("proto", temp[len(temp)-2], temp[len(temp)-1])
}

func ProtoPackage(module string) string {
	temp := strings.Split(module, "/")
	if len(temp) < 3 {
		return path.Base(module)
	}
	return strings.Join([]string{temp[len(temp)-2], temp[len(temp)-1]}, ".")
}
