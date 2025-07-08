package browser

import "go-drive/filesystem"

type doneLoadingMsg struct {
	filelist []filesystem.FileItem
	err      error
}

type syncStartMsg struct {
	index int
}

type syncDoneMsg struct {
	filelist []filesystem.FileItem
	index    int
	err      error
}
type syncFailedMsg struct {
	err   error
	index int
}

type trashStartMsg struct {
	index int
}

type trashDoneMsg struct {
	filelist []filesystem.FileItem
	index    int
	err      error
}
type trashFailedMsg struct {
	err   error
	index int
}
