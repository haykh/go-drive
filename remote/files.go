package remote

import (
	"go-drive/utils"

	"google.golang.org/api/drive/v3"
)

var RecognizedFormats = []string{
	"application/vnd.google-apps.folder",
	"application/vnd.google.colaboratory",
	"application/vnd.google-apps.document",
	"application/vnd.google-apps.spreadsheet",
	"application/vnd.google-apps.presentation",
	"application/pdf",
	"application/msword",
	"application/zip",
	"video/mp4",
	"image/png",
	"image/jpeg",
	"image/gif",
	"other",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
}

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

func (f File) GetName() string {
	return f.Name
}

func (f File) GetMimeType() string {
	return f.MimeType
}
