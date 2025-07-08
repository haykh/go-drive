package local

import (
	"go-drive/filesystem"
	"go-drive/utils"
	"time"
)

var _ filesystem.FileItem = &File{}
var _ filesystem.FileManager = &Manager{}

/* - - - - - - - - - -
 * Manager
 */
type Manager struct {
	Root string
}

func (m Manager) GetFileList(path string, debug_mode bool) ([]filesystem.FileItem, error) {
	if local_filelist, err := getFolderContent(m.Root, path); err != nil {
		return nil, utils.ToHumanReadableError(err, debug_mode)
	} else {
		wrappedFiles := make([]filesystem.FileItem, len(local_filelist))
		for i, f := range local_filelist {
			wrappedFiles[i] = f
		}
		return filesystem.Sorted(wrappedFiles), nil
	}
}

/* - - - - - - - - - -
 * File
 */
type File struct {
	FullPath     string
	RelativePath string
	Name         string
	Size         uint64
	ModifiedTime time.Time
	MimeType     string
	Md5Checksum  string
}

func (f File) IsDirectory() bool {
	return f.MimeType == "directory"
}

func (f File) IsPDF() bool {
	return f.MimeType == "application/pdf"
}

func (f File) IsUnrecognized() bool {
	return f.MimeType == "other"
}

func (f File) IsLocal() bool {
	return true
}

func (f File) IsRemote() bool {
	return false
}

func (f File) GetName() string {
	return f.Name
}

func (f File) GetSize() uint64 {
	return f.Size
}

func (f File) GetModifiedTime() time.Time {
	return f.ModifiedTime
}

func (f File) GetMimeType() string {
	return f.MimeType
}

func (f File) GetOwnedByMe() bool {
	return true
}
