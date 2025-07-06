package local

import (
	"go-drive/utils"
	"path/filepath"

	"os"

	"github.com/gabriel-vasile/mimetype"
)

func getFolderContent(root, path string) ([]*File, utils.APIError) {
	full_path := filepath.Join(root, path)
	if entries, err := os.ReadDir(full_path); err != nil {
		return nil, &utils.ReadDirFailed{OSError: err, Path: full_path}
	} else {
		files := []*File{}
		for _, entry := range entries {
			filesize := int64(0)
			modtime := ""
			mime_type := ""
			if !entry.IsDir() {
				if fileinfo, err := entry.Info(); err != nil {
					return nil, &utils.ReadFileInfoFailed{OSError: err, File: entry.Name(), Path: full_path}
				} else {
					filesize = fileinfo.Size()
					modtime = fileinfo.ModTime().GoString()
				}
				file := filepath.Join(full_path, entry.Name())
				if kind, err := mimetype.DetectFile(file); err != nil {
					return nil, &utils.MimeTypeFailed{OSError: err, File: file}
				} else {
					mime_type = kind.String()
				}
			} else {
				mime_type = "directory"
			}

			files = append(files, &File{
				FullPath:     full_path,
				RelativePath: path,
				Name:         entry.Name(),
				Size:         uint64(filesize),
				ModifiedTime: modtime,
				MimeType:     mime_type,
			})

		}
		return files, nil
	}
}
