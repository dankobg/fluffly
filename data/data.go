package data

import (
	"embed"
	"io/fs"
)

//go:embed public/*
var publicFS embed.FS

//go:embed templates/*
var templatesFS embed.FS

func PublicFS() (fs.FS, error) {
	return fs.Sub(publicFS, "public")
}

func MustPublicFS() fs.FS {
	subFS, err := PublicFS()
	if err != nil {
		panic("failed to get public fs subdir" + err.Error())
	}
	return subFS
}

func TemplatesFS() (fs.FS, error) {
	return fs.Sub(templatesFS, "templates")
}

func MustTemplatesFS() fs.FS {
	subFS, err := TemplatesFS()
	if err != nil {
		panic("failed to get public fs subdir" + err.Error())
	}
	return subFS
}
