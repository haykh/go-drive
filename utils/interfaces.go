package utils

import (
	"sort"
)

type FileManager interface {
	GetFileList(string, bool) ([]FileItem, error)
}

type FileItem interface {
	Stringize(string, string) string

	IsDirectory() bool
	IsPDF() bool

	GetName() string
	GetMimeType() string
}

func Sorted(f []FileItem) []FileItem {
	sort.Slice(f, func(i, j int) bool {
		if f[i].GetMimeType() == f[j].GetMimeType() {
			return f[i].GetName() < f[j].GetName()
		} else if f[i].IsDirectory() {
			return true
		} else if f[j].IsDirectory() {
			return false
		} else if f[i].IsPDF() {
			return true
		} else if f[j].IsPDF() {
			return false
		}
		if f[i].GetMimeType() == "other" {
			return true
		}
		if f[j].GetMimeType() == "other" {
			return false
		}
		return f[i].GetMimeType() < f[j].GetMimeType()
	})
	return f
}

func Stringize(f []FileItem, path, format string) []string {
	items := make([]string, len(f))
	for i, item := range f {
		items[i] = item.Stringize(path, format)
	}
	return items
}
