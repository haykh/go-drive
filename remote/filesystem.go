package remote

import (
	"fmt"
	"go-drive/utils"
	"strings"

	"github.com/charmbracelet/log"
	"google.golang.org/api/drive/v3"
)

func getFolder(srv *drive.Service, remote_path string) (*File, utils.APIError) {
	log.Debugf("getFolder: %s", remote_path)
	parentID := "root"
	parts := strings.Split(strings.Trim(remote_path, "/"), "/")
	var drive_file *drive.File

	for _, part := range parts {
		query := fmt.Sprintf("name = '%s' and mimeType = 'application/vnd.google-apps.folder' and '%s' in parents and trashed = false", part, parentID)
		if folder_obj, err := srv.Files.List().Q(query).Fields("files(id, name)").Do(); err != nil {
			return nil, &utils.QueryFailed{DriveError: err, Query: query}
		} else {
			if len(folder_obj.Files) == 0 {
				return nil, &utils.FolderNotFound{DriveError: err, Path: remote_path}
			} else {
				parentID = folder_obj.Files[0].Id
				drive_file = folder_obj.Files[0]
			}
		}
	}

	return &File{drive_file}, nil
}

func getFolderContent(srv *drive.Service, remote_path string) ([]*File, utils.APIError) {
	log.Debugf("getFolderContent: %s", remote_path)
	if remote_path == "" || remote_path == "/" {
		return queryFiles(srv, "'root' in parents and trashed = false")
	} else {
		if folder, err := getFolder(srv, remote_path); err != nil {
			return nil, err
		} else {
			return queryFiles(srv, fmt.Sprintf("'%s' in parents and trashed = false", folder.Id))
		}
	}
}

func queryFiles(srv *drive.Service, query string) ([]*File, utils.APIError) {
	if filelist, err := srv.Files.List().Q(query).Fields("files(id, name, mimeType, ownedByMe, modifiedTime, size)").Do(); err != nil {
		return nil, &utils.QueryFailed{DriveError: err, Query: query}
	} else {
		files := filelist.Files
		wrappedFiles := make([]*File, len(files))
		for i, f := range files {
			wrappedFiles[i] = &File{f}
		}
		return wrappedFiles, nil
	}
}

// func getFile(srv *drive.Service, file_name, remote_path string) (*File, utils.APIError) {
// 	log.Debugf("getFile: %s in %s", file_name, remote_path)
// 	if folder, err := getFolder(srv, remote_path); err != nil {
// 		return nil, err
// 	} else {
// 		return getFileInFolderId(srv, file_name, folder.Id, remote_path)
// 	}
// }

// func UploadFile(srv *drive.Service, file_path, remote_path string, mode utils.UploadMode) (*File, utils.APIError) {
// 	if file, err := os.Open(file_path); err != nil {
// 		return nil, &utils.OpenFileFailed{OSError: err, File: file_path}
// 	} else {
// 		defer file.Close()
// 		parts := strings.Split(strings.Trim(file_path, "/"), "/")
// 		file_name := parts[len(parts)-1]
// 		if remote_folder, err := getFolder(srv, remote_path); err != nil {
// 			return nil, err
// 		} else {
// 			if remote_file, err := getFileInFolderId(srv, file_name, remote_folder.Id, remote_path); err != nil {
// 				switch err.(type) {
// 				case *utils.FileNotFound:
// 					return createNewRemoteFile(srv,
// 						file_name,
// 						remote_path,
// 						&drive.File{
// 							Name:    file_name,
// 							Parents: []string{remote_folder.Id},
// 						},
// 						file,
// 					)
// 				default:
// 					return nil, err
// 				}
// 			} else {
// 				switch mode {
// 				case utils.RaiseIfDuplicate:
// 					return nil, &utils.DuplicateFile{File: file_name, Path: remote_path}
// 				case utils.SkipDuplicates:
// 					log.Warnf("File %s already exists in %s : skipping", file_name, remote_path)
// 					return remote_file, nil
// 				case utils.Overwrite:
// 					log.Warnf("File %s already exists in %s : overwriting", file_name, remote_path)
// 					return overwriteRemoteFile(
// 						srv,
// 						file_name,
// 						remote_file.Id,
// 						remote_path,
// 						&drive.File{},
// 						file,
// 					)
// 				default:
// 					return nil, &utils.WrongUploadMode{Mode: mode}
// 				}
// 			}
// 		}
// 	}
// }

// func getFileInFolderId(srv *drive.Service, file_name, folder_id, folder_name string) (*File, utils.APIError) {
// 	if filelist, err := queryFiles(srv, fmt.Sprintf("name = '%s' and '%s' in parents and trashed = false", file_name, folder_id)); err != nil {
// 		return nil, err
// 	} else {
// 		if len(filelist.Files) == 0 {
// 			return nil, &utils.FileNotFound{DriveError: nil, File: file_name, Path: folder_name}
// 		} else {
// 			return &File{filelist.Files[0]}, nil
// 		}
// 	}
// }

// func createNewRemoteFile(srv *drive.Service, file_name, remote_path string, metadata *drive.File, local_file *os.File) (*File, utils.APIError) {
// 	if newfile, err := srv.Files.Create(metadata).
// 		Media(local_file).
// 		Fields("id").
// 		Do(); err != nil {
// 		return nil, &utils.CreateFailed{DriveError: err, File: file_name, Path: remote_path}
// 	} else {
// 		return &File{newfile}, nil
// 	}
// }

// func overwriteRemoteFile(srv *drive.Service, file_name, remote_file_id, remote_path string, metadata *drive.File, local_file *os.File) (*File, utils.APIError) {
// 	if newfile, err := srv.Files.Update(remote_file_id, metadata).
// 		Media(local_file).
// 		Fields("id").
// 		Do(); err != nil {
// 		return nil, &utils.OverwriteFailed{DriveError: err, File: file_name, Path: remote_path}
// 	} else {
// 		return &File{newfile}, nil
// 	}
// }
