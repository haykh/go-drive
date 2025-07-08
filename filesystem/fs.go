package filesystem

import "time"

type FileManager interface {
	GetFileList(string, bool) ([]FileItem, error)

	Synchronize(FileItem, string, bool) error
	Trash(FileItem, string, bool) error
}

type FileItem interface {
	IsDirectory() bool
	IsPDF() bool
	IsUnrecognized() bool

	// IsLocal() bool
	// IsRemote() bool
	InSync() bool
	ShouldUpload() bool
	ShouldDownload() bool

	GetName() string
	GetSize() uint64
	GetMimeType() string
	GetModifiedTime() time.Time
	GetOwnedByMe() bool
}
