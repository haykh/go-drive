package remote

import (
	"fmt"
	"go-drive/utils"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"google.golang.org/api/drive/v3"
)

func getFolder(srv *drive.Service, remote_path string) (*File, utils.APIError) {
	log.Debugf("getFolder: %s", remote_path)
	parentID := "root"
	parts := strings.Split(strings.Trim(remote_path, "/"), "/")
	var drive_folder *drive.File

	for _, part := range parts {
		query := fmt.Sprintf("name = '%s' and mimeType = 'application/vnd.google-apps.folder' and '%s' in parents and trashed = false", part, parentID)
		if folder_obj, err := srv.Files.List().Q(query).Fields("files(id, name)").Do(); err != nil {
			return nil, &utils.QueryFailed{DriveError: err, Query: query}
		} else {
			if len(folder_obj.Files) == 0 {
				return nil, &utils.FolderNotFound{DriveError: err, Path: remote_path}
			} else {
				parentID = folder_obj.Files[0].Id
				drive_folder = folder_obj.Files[0]
			}
		}
	}
	return &File{drive_folder}, nil
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
	var all_files []*File
	page_token := ""
	for {
		request := srv.Files.List().Q(query).Fields("nextPageToken, files(id, name, mimeType, ownedByMe, modifiedTime, size, md5Checksum)")
		if page_token != "" {
			request = request.PageToken(page_token)
		}
		if filelist, err := request.Do(); err != nil {
			return nil, &utils.QueryFailed{DriveError: err, Query: query}
		} else {
			for _, f := range filelist.Files {
				all_files = append(all_files, &File{f})
			}
			if filelist.NextPageToken == "" {
				break
			}
			page_token = filelist.NextPageToken
		}
	}
	return all_files, nil
}

// func getFile(srv *drive.Service, file_name, remote_path string) (*File, utils.APIError) {
// 	log.Debugf("getFile: %s in %s", file_name, remote_path)
// 	if folder, err := getFolder(srv, remote_path); err != nil {
// 		return nil, err
// 	} else {
// 		return getFileInFolderId(srv, file_name, folder.Id, remote_path)
// 	}
// }

func getFileInFolderId(srv *drive.Service, file_name, folder_id, folder_name string) (*File, utils.APIError) {
	if filelist, err := queryFiles(srv, fmt.Sprintf("name = '%s' and '%s' in parents and trashed = false", file_name, folder_id)); err != nil {
		return nil, err
	} else {
		if len(filelist) == 0 {
			return nil, &utils.FileNotFound{DriveError: nil, File: file_name, Path: folder_name}
		} else {
			return filelist[0], nil
		}
	}
}

func createNewRemoteFile(srv *drive.Service, file_name, remote_path string, metadata *drive.File, local_file *os.File) (*File, utils.APIError) {
	if newfile, err := srv.Files.Create(metadata).
		Media(local_file).
		Fields("id").
		Do(); err != nil {
		return nil, &utils.CreateFailed{DriveError: err, File: file_name, Path: remote_path}
	} else {
		return &File{newfile}, nil
	}
}

func overwriteRemoteFile(srv *drive.Service, file_name, remote_file_id, remote_path string, metadata *drive.File, local_file *os.File) (*File, utils.APIError) {
	if newfile, err := srv.Files.Update(remote_file_id, metadata).
		Media(local_file).
		Fields("id").
		Do(); err != nil {
		return nil, &utils.OverwriteFailed{DriveError: err, File: file_name, Path: remote_path}
	} else {
		return &File{newfile}, nil
	}
}

func ensureFolderPath(srv *drive.Service, path string) (*drive.File, utils.APIError) {
	parts := strings.Split(path, "/")
	parentID := "root"
	var drive_folder *drive.File

	for _, part := range parts {
		if part == "" {
			continue
		}
		q := fmt.Sprintf("mimeType='application/vnd.google-apps.folder' and name='%s' and '%s' in parents and trashed=false", part, parentID)
		r, err := srv.Files.List().Q(q).Fields("files(id, name)").Do()
		if err != nil {
			return nil, &utils.QueryFailed{DriveError: err, Query: q}
		}

		if len(r.Files) > 0 {
			drive_folder = r.Files[0]
			parentID = r.Files[0].Id
		} else {
			f := &drive.File{
				Name:     part,
				MimeType: "application/vnd.google-apps.folder",
			}
			if parentID != "" {
				f.Parents = []string{parentID}
			}
			created, err := srv.Files.Create(f).Do()
			if err != nil {
				return nil, &utils.CreateFailed{
					DriveError: err,
					File:       part,
					Path:       path,
				}
			}
			drive_folder = created
			parentID = created.Id
		}
	}
	return drive_folder, nil
}

func moveFileById(srv *drive.Service, fileID, newParentID string) (*drive.File, utils.APIError) {
	file, err := srv.Files.Get(fileID).Fields("parents").Do()
	if err != nil {
		return nil, &utils.FileNotFound{
			DriveError: err,
			File:       fileID,
			Path:       "",
		}
	}

	oldParents := strings.Join(file.Parents, ",")

	new_file, err := srv.Files.Update(fileID, nil).
		AddParents(newParentID).
		RemoveParents(oldParents).
		Fields("id, parents").
		Do()

	if err != nil {
		return nil, &utils.CreateFailed{
			DriveError: err,
			File:       file.Name,
			Path:       newParentID,
		}
	}
	return new_file, nil
}
