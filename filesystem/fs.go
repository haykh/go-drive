package filesystem

type FileManager interface {
	GetFileList(string, bool) ([]FileItem, error)

	Synchronize(FileItem, bool) error
}

type FileItem interface {
	IsDirectory() bool
	IsPDF() bool
	IsUnrecognized() bool

	IsLocal() bool
	IsRemote() bool

	GetName() string
	GetSize() uint64
	GetMimeType() string
	GetModifiedTime() string
	GetOwnedByMe() bool
}
