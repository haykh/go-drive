package api

type FileMirror struct {
	RemoteId  string `json:"remote_id"`
	LocalPath string `json:"local_path"`
}

type SyncMap struct {
	Mirrors map[FileMirror]struct{}
}

// func (fs SyncMap) MarshalJSON() ([]byte, error) {
// 	var mirrors []FileMirror
// 	for m := range fs.Mirrors {
// 		mirrors = append(mirrors, m)
// 	}
// 	return json.Marshal(struct {
// 		Mirrors []FileMirror `json:"mirrors"`
// 	}{mirrors})
// }

// func (fs *SyncMap) UnmarshalJSON(data []byte) error {
// 	var aux struct {
// 		Mirrors []FileMirror `json:"mirrors"`
// 	}
// 	if err := json.Unmarshal(data, &aux); err != nil {
// 		return err
// 	}
// 	fs.Mirrors = make(map[FileMirror]struct{})
// 	for _, m := range aux.Mirrors {
// 		fs.Mirrors[m] = struct{}{}
// 	}
// 	return nil
// }

// func (f FileMirror) IsSynced() bool {
// 	return false
// }

// func NewSyncMap(srv *drive.Service, path string) (*SyncMap, utils.APIError) {
// 	fs := &SyncMap{
// 		Mirrors: map[FileMirror]struct{}{},
// 	}
// 	if itemlist, err := remote.GetFileList(srv, path); err != nil {
// 		return nil, err
// 	} else {
// 		for _, item := range itemlist.Files {
// 			if item.MimeType == "application/vnd.google-apps.folder" {
// 				if subfs, err := NewSyncMap(srv, filepath.Join(path, item.Name)); err != nil {
// 					return nil, err
// 				} else {
// 					for mirror := range subfs.Mirrors {
// 						fs.Mirrors[mirror] = struct{}{}
// 					}
// 				}
// 			} else {
// 				fs.Mirrors[FileMirror{
// 					RemoteId:  item.Id,
// 					LocalPath: filepath.Join("/", path, item.Name),
// 				}] = struct{}{}
// 			}
// 		}
// 	}
// 	return fs, nil
// }

// func (fs SyncMap) ToJson(json_path string) utils.APIError {
// 	if file, err := os.Create(json_path); err != nil {
// 		return &utils.CreateFileFailed{OSError: err, File: json_path}
// 	} else {
// 		defer file.Close()
// 		if data, err := json.Marshal(fs); err != nil {
// 			return &utils.JSONMarshalFailed{OSError: err, Name: "filesystem"}
// 		} else if _, err := file.Write(data); err != nil {
// 			return &utils.WriteFileFailed{OSError: err, File: json_path}
// 		}
// 		return nil
// 	}
// }

// func (fs *SyncMap) FromJson(json_path string) utils.APIError {
// 	if file, err := os.Open(json_path); err != nil {
// 		return &utils.OpenFileFailed{OSError: err, File: json_path}
// 	} else {
// 		defer file.Close()
// 		if data, err := os.ReadFile(json_path); err != nil {
// 			return &utils.ReadFileFailed{OSError: err, File: json_path}
// 		} else if err := json.Unmarshal(data, fs); err != nil {
// 			return &utils.JSONUnmarshalFailed{OSError: err, File: json_path}
// 		}
// 		return nil
// 	}
// }

// func GetSyncMap(srv *drive.Service, json_path string) (*SyncMap, utils.APIError) {
// 	fs := &SyncMap{
// 		Mirrors: map[FileMirror]struct{}{},
// 	}
// 	if err := fs.FromJson(json_path); err != nil {
// 		if new_fs, err := NewSyncMap(srv, ""); err != nil {
// 			return nil, err
// 		} else if err := new_fs.ToJson(json_path); err != nil {
// 			return nil, err
// 		} else {
// 			fs = new_fs
// 		}
// 	}
// 	return fs, nil
// }
