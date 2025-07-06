package filesystem

import "sort"

func Sorted(f []FileItem) []FileItem {
	sort.Slice(f, func(i, j int) bool {
		if f[i].GetMimeType() == f[j].GetMimeType() {
			return f[i].GetName() < f[j].GetName()
		} else if f[i].IsDirectory() || f[j].IsDirectory() {
			return f[i].IsDirectory()
		} else if f[i].IsPDF() || f[j].IsPDF() {
			return f[i].IsPDF()
		} else if f[i].GetMimeType() == "other" || f[j].GetMimeType() == "other" {
			return f[j].GetMimeType() == "other"
		}
		return f[i].GetMimeType() < f[j].GetMimeType()
	})
	return f
}
