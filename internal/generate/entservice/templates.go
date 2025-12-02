package entservice

import (
	"embed"
)

//go:embed templates
var fs embed.FS

type output struct {
	filename  string
	overwrite bool
}

func out(filename string, overwrite ...bool) *output {
	return &output{filename: filename, overwrite: len(overwrite) > 0 && overwrite[0]}
}
