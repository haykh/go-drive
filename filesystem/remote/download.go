package remote

import (
	"go-drive/utils"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"google.golang.org/api/drive/v3"
)

func DownloadFile(srv *drive.Service, local_path, remote_filepath string, mode utils.SyncMode) (*os.File, utils.APIError) {
	log.Debugf("DownloadFile: %s from %s", local_path, remote_filepath)

	// get remote file name and folder path
	parts := strings.Split(strings.Trim(remote_filepath, "/"), "/")
	file_name := parts[len(parts)-1]
	remote_path := strings.Join(parts[:len(parts)-1], "/")

	// get remote folder Id
	remote_folder_id := "root"
	if remote_path != "" && remote_path != "/" {
		if remote_folder, err := getFolder(srv, remote_path); err != nil {
			return nil, err
		} else {
			remote_folder_id = remote_folder.Id
		}
	}

	// get remote file id
	if remote_file, err := getFileInFolderId(srv, file_name, remote_folder_id, remote_path); err != nil {
		return nil, err
	} else {
		// check if file already exists locally
		local_filepath := filepath.Join(local_path, file_name)
		if file, err := os.Open(local_filepath); err == nil {
			file.Close()
			// if file exists, handle according to mode
			switch mode {
			case utils.RaiseIfDuplicate:
				return nil, &utils.DuplicateFile{File: file_name, Path: remote_path}
			case utils.SkipDuplicates:
				log.Warnf("File %s already exists in %s : skipping", file_name, remote_path)
				return nil, nil
			case utils.Overwrite:
				log.Warnf("File %s already exists in %s : overwriting", file_name, remote_path)
				if err := os.Remove(local_filepath); err != nil {
					return nil, &utils.RemoveFileFailed{OSError: err, File: local_filepath}
				}
			}
		}

		// download the file
		if res, err := srv.Files.Get(remote_file.Id).Download(); err != nil {
			return nil, &utils.DownloadFailed{DriveError: err, File: remote_file.Name}
		} else {
			defer res.Body.Close()

			// create local directory
			if err := os.MkdirAll(local_path, 0755); err != nil {
				return nil, &utils.CreateDirFailed{OSError: err, Dir: local_path}
			}

			// create local file
			if local_file, err := os.Create(local_filepath); err != nil {
				return nil, &utils.CreateFileFailed{OSError: err, File: local_filepath}
			} else {
				defer local_file.Close()

				// copy content from remote to local file
				if _, err = io.Copy(local_file, res.Body); err != nil {
					return nil, &utils.CopyFileFailed{OSError: err, File: local_filepath}
				}
				if err := os.Chtimes(local_filepath, time.Now().UTC(), remote_file.GetModifiedTime()); err != nil {
					return nil, &utils.ChtimeFailed{OSError: err, File: local_filepath}
				}
				return local_file, nil
			}
		}
	}
}
