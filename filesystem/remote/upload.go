package remote

import (
	"go-drive/utils"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"google.golang.org/api/drive/v3"
)

func UploadFile(srv *drive.Service, local_filepath, remote_path string, mode utils.SyncMode) (*File, utils.APIError) {
	log.Debugf("UploadFile: %s to %s", local_filepath, remote_path)

	// open the local file
	if file, err := os.Open(local_filepath); err != nil {
		return nil, &utils.OpenFileFailed{OSError: err, File: local_filepath}
	} else {
		defer file.Close()

		// get local file name and folder path
		parts := strings.Split(strings.Trim(local_filepath, "/"), "/")
		file_name := parts[len(parts)-1]

		// get remote folder Id
		remote_folder_id := "root"
		if remote_path != "" && remote_path != "/" {
			if remote_folder, err := getFolder(srv, remote_path); err != nil {
				return nil, err
			} else {
				remote_folder_id = remote_folder.Id
			}
		}

		// check if file already exists in remote folder
		if remote_file, err := getFileInFolderId(srv, file_name, remote_folder_id, remote_path); err != nil {
			// if file does not exist, create a new remote file
			switch err.(type) {
			case *utils.FileNotFound:
				return createNewRemoteFile(srv,
					file_name,
					remote_path,
					&drive.File{
						Name:    file_name,
						Parents: []string{remote_folder_id},
					},
					file,
				)
			default:
				return nil, err
			}
		} else {
			// if file exists, handle according to mode
			switch mode {
			case utils.RaiseIfDuplicate:
				return nil, &utils.DuplicateFile{File: file_name, Path: remote_path}
			case utils.SkipDuplicates:
				log.Warnf("File %s already exists in %s : skipping", file_name, remote_path)
				return remote_file, nil
			case utils.Overwrite:
				log.Warnf("File %s already exists in %s : overwriting", file_name, remote_path)
				return overwriteRemoteFile(
					srv,
					file_name,
					remote_file.Id,
					remote_path,
					&drive.File{},
					file,
				)
			default:
				return nil, &utils.WrongSyncMode{Mode: mode}
			}
		}
	}
}
