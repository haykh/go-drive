package local

import (
	"go-drive/utils"
	"path/filepath"
)

var RecognizedFormats = []string{}

var _ utils.FileItem = &File{}
var _ utils.FileManager = &Manager{}

/* - - - - - - - - - -
 * Manager
 */
type Manager struct {
	Root string
}

func (m Manager) GetFileList(path string, debug_mode bool) ([]utils.FileItem, error) {
	if local_filelist, err := getFolderContent(filepath.Join(m.Root, path)); err != nil {
		return nil, utils.ToHumanReadableError(err, debug_mode)
	} else {
		wrappedFiles := make([]utils.FileItem, len(local_filelist))
		for i, f := range local_filelist {
			wrappedFiles[i] = f
		}
		return utils.Sorted(wrappedFiles), nil
	}
}

/* - - - - - - - - - -
 * File
 */
type File struct {
	Path         string
	Name         string
	Size         int64
	ModifiedTime string
	MimeType     string
}

func (f File) IsDirectory() bool {
	return f.MimeType == "directory"
}

func (f File) IsPDF() bool {
	return f.MimeType == "application/pdf"
}

func (f File) GetName() string {
	return f.Name
}

func (f File) GetMimeType() string {
	return f.MimeType
}
