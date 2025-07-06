package api

import (
	"bytes"
	"go-drive/components/browser"
	"go-drive/components/spinner"
	"go-drive/filesystem"
	"go-drive/filesystem/dual"
	"go-drive/filesystem/local"
	"go-drive/filesystem/remote"
	"strings"

	"github.com/charmbracelet/log"
	"google.golang.org/api/drive/v3"
)

func DualLs(srv *drive.Service, local_root, dir string, debug_mode bool) error {
	mgr := dual.DualManager{
		LocalManager:  &local.Manager{Root: local_root},
		RemoteManager: &remote.Manager{Srv: srv},
	}
	if content, err := spinner.RunWithSpinner(
		func() (any, error) {
			return mgr.GetFileList(dir, debug_mode)
		},
		"loading",
		"unable to get local or remote content",
		"",
		debug_mode,
	); err != nil {
		return err
	} else {
		itemlist, _ := content.([]filesystem.FileItem)
		log.Print(strings.Join(filesystem.StringizeAll(itemlist, dir), "\n"))
		return nil
	}
}

// func RemoteFileBrowser(srv *drive.Service, dir string, debug_mode bool) error {
// 	return browser.FileBrowser(&remote.Manager{Srv: srv}, dir, debug_mode)
// }

// func LocalFileBrowser(local_root, dir string, debug_mode bool) error {
// 	return browser.FileBrowser(&local.Manager{Root: local_root}, dir, debug_mode)
// }

func DualFileBrowser(srv *drive.Service, local_root, dir string, debug_mode bool, debugBuffer *bytes.Buffer) error {
	return browser.FileBrowser(
		&dual.DualManager{
			LocalManager:  &local.Manager{Root: local_root},
			RemoteManager: &remote.Manager{Srv: srv},
		},
		dir,
		debug_mode,
		debugBuffer,
	)
}

// func getSyncmap(srv *drive.Service, syncmap_path string, debug_mode bool) (*SyncMap, error) {
// 	if content, err := ui.RunWithSpinner(
// 		func() (any, error) {
// 			if fs, err := GetSyncMap(srv, syncmap_path); err != nil {
// 				return nil, utils.ToHumanReadableError(err, debug_mode)
// 			} else {
// 				return fs, nil
// 			}
// 		},
// 		"building local filesystem",
// 		"unable to build local filesystem",
// 		"",
// 		debug_mode,
// 	); err != nil {
// 		return nil, err
// 	} else {
// 		return content.(*SyncMap), nil
// 	}
// }
