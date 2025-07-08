package dual

import (
	"fmt"
	"go-drive/filesystem"
)

func (m DualManager) Trash(file filesystem.FileItem, path string, debug_mode bool) error {
	if df, ok := file.(DualFile); ok {
		var err_remote, err_local error
		if df.RemoteFile != nil {
			err_remote = m.RemoteManager.Trash(df.RemoteFile, path, debug_mode)
		}
		if df.LocalFile != nil {
			err_local = m.LocalManager.Trash(df.LocalFile, path, debug_mode)
		}
		if err_remote != nil || err_local != nil {
			return fmt.Errorf("failed to trash file %s: remote error: %v, local error: %v", file.GetName(), err_remote, err_local)
		}
		return nil
	} else {
		return fmt.Errorf("file is not a dual file: %T", file)
	}
}
