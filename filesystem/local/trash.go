package local

import (
	"fmt"
	"go-drive/filesystem"
	"go-drive/utils"
	"os"
	"path/filepath"
)

func (m Manager) Trash(file filesystem.FileItem, path string, debug_mode bool) error {
	trash_dir, err := ensureFolderPath(m.Root, filepath.Join(".Trash", path))
	if err != nil {
		return utils.ToHumanReadableError(err, debug_mode)
	}
	if lf, ok := file.(*File); ok {
		oldfile_path := filepath.Join(m.Root, path, lf.Name)
		newfile_path := filepath.Join(trash_dir, lf.Name)
		return os.Rename(oldfile_path, newfile_path)
	} else {
		return fmt.Errorf("file is not a local file: %T", file)
	}
}
