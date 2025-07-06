package browser

import "go-drive/filesystem"

type doneLoadingMsg struct {
	filelist []filesystem.FileItem
	err      error
}

type doneSyncMsg struct {
	filelist []filesystem.FileItem
	index    int
	err      error
}

type syncFileMsg struct {
	index int
}

type errorMsg struct {
	err error
}
