package api

import (
	"fmt"
	"go-drive/browser"
	"go-drive/local"
	"go-drive/remote"
	"go-drive/ui"
	"go-drive/utils"
	"strings"

	"github.com/charmbracelet/log"
	"google.golang.org/api/drive/v3"
)

func RemoteLs(srv *drive.Service, dir string, debug_mode bool) error {
	mgr := remote.Manager{Srv: srv}
	if content, err := ui.RunWithSpinner(
		func() (any, error) {
			return mgr.GetFileList(dir, debug_mode)
		},
		"loading",
		"unable to get directory content",
		"",
		debug_mode,
	); err != nil {
		return err
	} else {
		itemlist, _ := content.([]utils.FileItem)
		log.Print(strings.Join(utils.StringizeAll(itemlist, dir), "\n"))
		return nil
	}
}

func RemoteFileBrowser(srv *drive.Service, syncmap_path, dir string, debug_mode bool) error {
	// if _, err := getSyncmap(srv, syncmap_path, debug_mode); err != nil {
	// 	return err
	// } else {
	return browser.FileBrowser(&remote.Manager{Srv: srv}, dir, debug_mode)
	// }
}

func LocalFileBrowser(syncmap_path, root, dir string, debug_mode bool) error {
	fmt.Printf("DIRECTORY IS %s\n", dir)
	return browser.FileBrowser(&local.Manager{Root: root}, dir, debug_mode)
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
