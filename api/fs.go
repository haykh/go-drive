package api

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"google.golang.org/api/drive/v3"
)

type UploadMode int

const (
	RaiseIfDuplicate UploadMode = iota
	SkipDuplicates
	Overwrite
)

func queryFiles(srv *drive.Service, query string) (*drive.FileList, APIError) {
	if filelist, err := srv.Files.List().Q(query).Fields("files(id, name, mimeType, ownedByMe, size)").Do(); err != nil {
		return nil, &QueryFailed{err, query}
	} else {
		return filelist, nil
	}
}

func getFileInFolderId(srv *drive.Service, file_name, folder_id, folder_name string) (*drive.File, APIError) {
	if filelist, err := queryFiles(srv, fmt.Sprintf("name = '%s' and '%s' in parents and trashed = false", file_name, folder_id)); err != nil {
		return nil, err
	} else {
		if len(filelist.Files) == 0 {
			return nil, &FileNotFound{nil, file_name, folder_name}
		} else {
			return filelist.Files[0], nil
		}
	}
}

func createNewRemoteFile(srv *drive.Service, file_name, remote_path string, metadata *drive.File, local_file *os.File) (*drive.File, APIError) {
	if newfile, err := srv.Files.Create(metadata).
		Media(local_file).
		Fields("id").
		Do(); err != nil {
		return nil, &CreateFailed{err, file_name, remote_path}
	} else {
		return newfile, nil
	}
}

func overwriteRemoteFile(srv *drive.Service, file_name, remote_file_id, remote_path string, metadata *drive.File, local_file *os.File) (*drive.File, APIError) {
	if newfile, err := srv.Files.Update(remote_file_id, metadata).
		Media(local_file).
		Fields("id").
		Do(); err != nil {
		return nil, &OverwriteFailed{err, file_name, remote_path}
	} else {
		return newfile, nil
	}
}

func GetFolder(srv *drive.Service, remote_path string) (*drive.File, APIError) {
	log.Debugf("GetFolder: %s", remote_path)
	parentID := "root"
	parts := strings.Split(strings.Trim(remote_path, "/"), "/")
	var drive_file *drive.File

	for _, part := range parts {
		query := fmt.Sprintf("name = '%s' and mimeType = 'application/vnd.google-apps.folder' and '%s' in parents and trashed = false", part, parentID)
		if folder_obj, err := srv.Files.List().Q(query).Fields("files(id, name)").Do(); err != nil {
			return nil, &QueryFailed{err, query}
		} else {
			if len(folder_obj.Files) == 0 {
				return nil, &FolderNotFound{err, remote_path}
			} else {
				parentID = folder_obj.Files[0].Id
				drive_file = folder_obj.Files[0]
			}
		}
	}

	return drive_file, nil
}

func GetFolderContent(srv *drive.Service, remote_path string) (*drive.FileList, APIError) {
	log.Debugf("GetFolderContent: %s", remote_path)
	if remote_path == "" || remote_path == "/" {
		return queryFiles(srv, "'root' in parents and trashed = false")
	} else {
		if folder, err := GetFolder(srv, remote_path); err != nil {
			return nil, err
		} else {
			return queryFiles(srv, fmt.Sprintf("'%s' in parents and trashed = false", folder.Id))
		}
	}
}

func GetFile(srv *drive.Service, file_name, remote_path string) (*drive.File, APIError) {
	log.Debugf("GetFile: %s in %s", file_name, remote_path)
	if folder, err := GetFolder(srv, remote_path); err != nil {
		return nil, err
	} else {
		return getFileInFolderId(srv, file_name, folder.Id, remote_path)
	}
}

func UploadFile(srv *drive.Service, file_path, remote_path string, mode UploadMode) (*drive.File, APIError) {
	if file, err := os.Open(file_path); err != nil {
		return nil, &OpenFileFailed{err, file_path}
	} else {
		defer file.Close()
		parts := strings.Split(strings.Trim(file_path, "/"), "/")
		file_name := parts[len(parts)-1]
		if remote_folder, err := GetFolder(srv, remote_path); err != nil {
			return nil, err
		} else {
			if remote_file, err := getFileInFolderId(srv, file_name, remote_folder.Id, remote_path); err != nil {
				switch err.(type) {
				case *FileNotFound:
					return createNewRemoteFile(srv,
						file_name,
						remote_path,
						&drive.File{
							Name:    file_name,
							Parents: []string{remote_folder.Id},
						},
						file,
					)
				default:
					return nil, err
				}
			} else {
				switch mode {
				case RaiseIfDuplicate:
					return nil, &DuplicateFile{file_name, remote_path}
				case SkipDuplicates:
					log.Warnf("File %s already exists in %s : skipping", file_name, remote_path)
					return remote_file, nil
				case Overwrite:
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
					return nil, &WrongUploadMode{mode}
				}
			}
		}
	}
}
