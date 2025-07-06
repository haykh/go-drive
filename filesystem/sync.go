package filesystem

import (
	"fmt"
	"go-drive/filesystem/remote"
	"go-drive/utils"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func (m DualManager) Synchronize(file utils.FileItem, debug_mode bool) error {
	log.Debugf("Synchronizing file: %s", file.GetName())
	if file.IsLocal() && file.IsRemote() {
		return nil
	} else if file.IsLocal() {
		if f, ok := file.(DualFile); !ok {
			return fmt.Errorf("cannot synchronize: file is not a DualFile")
		} else {
			local_path := filepath.Join(f.LocalFile.FullPath, f.LocalFile.Name)
			remote_path := f.LocalFile.RelativePath
			if _, err := remote.UploadFile(
				m.RemoteManager.Srv,
				local_path,
				remote_path,
				utils.RaiseIfDuplicate,
			); err != nil {
				return utils.ToHumanReadableError(err, debug_mode)
			}
			return nil
		}
	} else if file.IsRemote() {
		return nil
	}
	return fmt.Errorf("cannot synchronize: file is neither local nor remote")
}
