package remote

import (
	"go-drive/ui"
	"go-drive/utils"

	"google.golang.org/api/drive/v3"
)

var _ utils.FileItem = &File{}
var _ utils.FileManager = &Manager{}

/* - - - - - - - - - -
 * Manager
 */
type Manager struct {
	Srv *drive.Service
}

func (m Manager) GetFileList(path string, debug_mode bool) ([]utils.FileItem, error) {
	if remote_filelist, err := getFolderContent(m.Srv, path); err != nil {
		return nil, utils.ToHumanReadableError(err, debug_mode)
	} else {
		wrappedFiles := make([]utils.FileItem, len(remote_filelist))
		for i, f := range remote_filelist {
			wrappedFiles[i] = f
		}
		return utils.Sorted(wrappedFiles), nil
	}
}

/* - - - - - - - - - -
 * File
 */
type File struct {
	*drive.File
}

func (f File) IsDirectory() bool {
	return f.MimeType == "application/vnd.google-apps.folder"
}

func (f File) IsPDF() bool {
	return f.MimeType == "application/pdf"
}

func (f File) IsUnrecognized() bool {
	if _, ok := ui.MimeIcons[f.MimeType]; !ok {
		return true
	} else {
		return f.MimeType == "other"
	}
}

func (f File) IsLocal() bool {
	return false
}

func (f File) IsRemote() bool {
	return true
}

func (f File) GetName() string {
	return f.Name
}

func (f File) GetSize() uint64 {
	return uint64(f.Size)
}

func (f File) GetMimeType() string {
	return f.MimeType
}

func (f File) GetModifiedTime() string {
	return f.ModifiedTime
}

func (f File) GetOwnedByMe() bool {
	return f.OwnedByMe
}
