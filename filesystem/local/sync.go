package local

import "go-drive/filesystem"

func (f File) InSync() bool {
	return true
}

func (f File) ShouldUpload() bool {
	return false
}

func (f File) ShouldDownload() bool {
	return false
}

func (m Manager) Synchronize(filesystem.FileItem, string, bool) error {
	return nil
}
