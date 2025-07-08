package api

import (
	"bytes"
	"fmt"
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

type LSMode int

const (
	LSDual LSMode = iota
	LSRemote
	LSLocal
)

func ListFiles(srv *drive.Service, ls_mode LSMode, local_root, dir string, debug_mode bool) error {
	var loader func() (any, error)
	switch ls_mode {
	case LSRemote:
		loader = func() (any, error) {
			return remote.Manager{Srv: srv}.GetFileList(dir, debug_mode)
		}
	case LSLocal:
		loader = func() (any, error) {
			return local.Manager{Root: local_root}.GetFileList(dir, debug_mode)
		}
	case LSDual:
		loader = func() (any, error) {
			mgr := dual.DualManager{
				LocalManager:  &local.Manager{Root: local_root},
				RemoteManager: &remote.Manager{Srv: srv},
			}
			return mgr.GetFileList(dir, debug_mode)
		}
	default:
		return fmt.Errorf("invalid LSMode: %v", ls_mode)
	}
	if content, err := spinner.RunWithSpinner(
		loader,
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

// func RemoteLs(srv *drive.Service, local_root, dir string, debug_mode bool) error {
// 	if content, err := spinner.RunWithSpinner(
// 		,
// 		"loading",
// 		"unable to get remote content",
// 		"",
// 		debug_mode,
// 	); err != nil {
// 		return err
// 	} else {
// 		itemlist, _ := content.([]filesystem.FileItem)
// 		log.Print(strings.Join(filesystem.StringizeAll(itemlist, dir), "\n"))
// 		return nil
// 	}
// }

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
