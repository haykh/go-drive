package local

import (
	"crypto/md5"
	"encoding/hex"
	"go-drive/utils"
	"io"
	"path/filepath"
	"time"

	"os"

	"github.com/gabriel-vasile/mimetype"
)

func md5sum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func getFolderContent(root, path string) ([]*File, utils.APIError) {
	full_path := filepath.Join(root, path)
	if entries, err := os.ReadDir(full_path); err != nil {
		return nil, &utils.ReadDirFailed{OSError: err, Path: full_path}
	} else {
		files := []*File{}
		for _, entry := range entries {
			filesize := int64(0)
			modtime := time.Time{}
			mime_type := ""
			md5_checksum := ""
			if !entry.IsDir() {
				if fileinfo, err := entry.Info(); err != nil {
					return nil, &utils.ReadFileInfoFailed{OSError: err, File: entry.Name(), Path: full_path}
				} else {
					filesize = fileinfo.Size()
					modtime = fileinfo.ModTime().UTC()
				}
				file := filepath.Join(full_path, entry.Name())
				if kind, err := mimetype.DetectFile(file); err != nil {
					return nil, &utils.MimeTypeFailed{OSError: err, File: file}
				} else {
					mime_type = kind.String()
				}
				if checksum, err := md5sum(file); err != nil {
					return nil, &utils.Md5Failed{OSError: err, File: file}
				} else {
					md5_checksum = checksum
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
				Md5Checksum:  md5_checksum,
			})

		}
		return files, nil
	}
}

func ensureFolderPath(root, path string) (string, utils.APIError) {
	full_path := filepath.Join(root, path)
	if err := os.MkdirAll(full_path, 0755); err != nil {
		return "", &utils.CreateDirFailed{OSError: err, Dir: full_path}
	}
	return full_path, nil
}
