package remote

import (
	"fmt"
	"go-drive/filesystem"
	"go-drive/utils"
	"strings"
)

func (m Manager) Trash(file filesystem.FileItem, path string, debug_mode bool) error {
	trash_path := strings.Join([]string{".Trash", path}, "/")
	trash_dir, err := ensureFolderPath(m.Srv, trash_path)
	if err != nil {
		return utils.ToHumanReadableError(err, debug_mode)
	}
	if rf, ok := file.(*File); ok {
		if rf.Trashed {
			return nil
		}
		if _, err := moveFileById(m.Srv, rf.Id, trash_dir.Id); err != nil {
			return utils.ToHumanReadableError(err, debug_mode)
		}
		return nil
	} else {
		return fmt.Errorf("file is not a remote file: %T", file)
	}
}
