package dual

import (
	"fmt"
	"go-drive/filesystem"
	"go-drive/filesystem/remote"
	"go-drive/utils"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func (f DualFile) InSync() bool {
	if f.RemoteFile == nil || f.LocalFile == nil {
		return false
	}
	return f.RemoteFile.Md5Checksum == f.LocalFile.Md5Checksum
}

func (f DualFile) ShouldUpload() bool {
	if f.InSync() {
		return false
	} else {
		if f.RemoteFile == nil {
			return true
		} else if f.LocalFile == nil {
			return false
		} else {
			return f.LocalFile.GetModifiedTime().After(f.RemoteFile.GetModifiedTime())
		}
	}
}

func (f DualFile) ShouldDownload() bool {
	if f.InSync() {
		return false
	} else {
		return !f.ShouldUpload()
	}
}

func (m DualManager) Synchronize(file filesystem.FileItem, relative_path string, debug_mode bool) error {
	log.Debugf("Synchronizing file: %s", file.GetName())
	if file.InSync() {
		return nil
	} else if file.ShouldUpload() {
		if f, ok := file.(DualFile); !ok {
			return fmt.Errorf("cannot synchronize: file is not a DualFile")
		} else {
			local_filepath := filepath.Join(f.LocalFile.FullPath, f.LocalFile.Name)
			if _, err := remote.UploadFile(
				m.RemoteManager.Srv,
				local_filepath,
				relative_path,
				utils.Overwrite,
			); err != nil {
				return utils.ToHumanReadableError(err, debug_mode)
			}
			return nil
		}
	} else if file.ShouldDownload() {
		if f, ok := file.(DualFile); !ok {
			return fmt.Errorf("cannot synchronize: file is not a DualFile")
		} else {
			remote_filepath := filepath.Join(relative_path, f.RemoteFile.Name)
			local_path := filepath.Join(m.LocalManager.Root, relative_path)
			if _, err := remote.DownloadFile(
				m.RemoteManager.Srv,
				local_path,
				remote_filepath,
				utils.Overwrite,
			); err != nil {
				return utils.ToHumanReadableError(err, debug_mode)
			}
			return nil
		}
	}
	return fmt.Errorf("cannot synchronize: file is neither local nor remote")
}
