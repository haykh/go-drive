package filesystem

import (
	"go-drive/filesystem/local"
	"go-drive/filesystem/remote"
	"go-drive/utils"
)

var _ utils.FileItem = &DualFile{}
var _ utils.FileManager = &DualManager{}

/* - - - - - - - - - -
 * Manager
 */
type DualManager struct {
	RemoteManager *remote.Manager
	LocalManager  *local.Manager
}

func (m DualManager) GetFileList(path string, debug_mode bool) ([]utils.FileItem, error) {
	remote_filelist, err := m.RemoteManager.GetFileList(path, debug_mode)
	if err != nil {
		return nil, err
	}
	local_filelist, err := m.LocalManager.GetFileList(path, debug_mode)
	if err != nil {
		return nil, err
	}
	dual_filelist := []utils.FileItem{}
	for _, remote_file := range remote_filelist {
		dual_file := DualFile{
			RemoteFile: remote_file.(*remote.File),
			LocalFile:  nil,
		}
		for lidx, local_file := range local_filelist {
			if (remote_file.GetName() == local_file.GetName()) && (remote_file.IsDirectory() == local_file.IsDirectory()) {
				dual_file.LocalFile = local_file.(*local.File)
				local_filelist = append(local_filelist[:lidx], local_filelist[lidx+1:]...)
				break
			}
		}
		dual_filelist = append(dual_filelist, dual_file)
	}
	for _, local_file := range local_filelist {
		dual_file := DualFile{
			RemoteFile: nil,
			LocalFile:  local_file.(*local.File),
		}
		dual_filelist = append(dual_filelist, dual_file)
	}
	return utils.Sorted(dual_filelist), nil
}

/* - - - - - - - - - -
 * File
 */
type DualFile struct {
	RemoteFile *remote.File
	LocalFile  *local.File
}

func (f DualFile) agree(remoteFn, localFn func() bool) bool {
	switch {
	case f.RemoteFile != nil && f.LocalFile != nil:
		rv := remoteFn()
		lv := localFn()
		if rv != lv {
			panic("remote and local files disagree")
		}
		return rv
	case f.RemoteFile != nil:
		return remoteFn()
	case f.LocalFile != nil:
		return localFn()
	default:
		panic("both remote and local are nil")
	}
}

func (f DualFile) IsDirectory() bool {
	return f.agree(
		func() bool { return f.RemoteFile.IsDirectory() },
		func() bool { return f.LocalFile.IsDirectory() },
	)
}

func (f DualFile) IsPDF() bool {
	return f.agree(
		func() bool { return f.RemoteFile.IsPDF() },
		func() bool { return f.LocalFile.IsPDF() },
	)
}

func (f DualFile) IsUnrecognized() bool {
	return f.agree(
		func() bool { return f.RemoteFile.IsUnrecognized() },
		func() bool { return f.LocalFile.IsUnrecognized() },
	)
}

func (f DualFile) IsLocal() bool {
	return f.LocalFile != nil
}

func (f DualFile) IsRemote() bool {
	return f.RemoteFile != nil
}

func (f DualFile) GetName() string {
	switch {
	case f.RemoteFile != nil && f.LocalFile != nil:
		if f.RemoteFile.GetName() != f.LocalFile.GetName() {
			panic("remote and local files disagree on name")
		}
		return f.RemoteFile.GetName()
	case f.RemoteFile != nil:
		return f.RemoteFile.GetName()
	case f.LocalFile != nil:
		return f.LocalFile.GetName()
	default:
		panic("both remote and local files are nil")
	}
}

func (f DualFile) GetSize() uint64 {
	if f.RemoteFile != nil {
		return f.RemoteFile.GetSize()
	} else if f.LocalFile != nil {
		return f.LocalFile.GetSize()
	} else {
		panic("both remote and local files are nil")
	}
}

func (f DualFile) GetMimeType() string {
	switch {
	case f.RemoteFile != nil:
		return f.RemoteFile.GetMimeType()
	case f.LocalFile != nil:
		return f.LocalFile.GetMimeType()
	default:
		panic("both remote and local files are nil")
	}
}

func (f DualFile) GetModifiedTime() string {
	switch {
	case f.RemoteFile != nil:
		return f.RemoteFile.GetModifiedTime()
	case f.LocalFile != nil:
		return f.LocalFile.GetModifiedTime()
	default:
		panic("both remote and local files are nil")
	}
}

func (f DualFile) GetOwnedByMe() bool {
	if f.RemoteFile != nil {
		return f.RemoteFile.GetOwnedByMe()
	} else {
		return true
	}
}
