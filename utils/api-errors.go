package utils

import (
	"fmt"

	"github.com/charmbracelet/log"
)

type APIError interface {
	HumanReadableError(bool) error
}

func ToHumanReadableError(err APIError, debug bool) error {
	if err != nil && debug {
		return err.HumanReadableError(debug)
	}
	return nil
}

// Errors from the Google Drive API
type GoogleDriveError struct {
	DriveError error
}

func (e *GoogleDriveError) HumanReadableError(debug bool) error {
	log.Errorf("GoogleDriveError: %v", e.DriveError)
	if debug {
		return e.DriveError
	} else {
		return nil
	}
}

type ParseCredentialsFailed struct {
	DriveError error
}

func (e *ParseCredentialsFailed) HumanReadableError(debug bool) error {
	log.Errorf("ParseCredentialsFailed: %v", e.DriveError)
	if debug {
		return e.DriveError
	} else {
		return nil
	}
}

type DownloadFailed struct {
	DriveError error
	File       string
}

func (e *DownloadFailed) HumanReadableError(debug bool) error {
	log.Errorf("DownloadFailed: %s", e.File)
	if debug {
		return e.DriveError
	} else {
		return nil
	}
}

type AuthTokenFailed struct {
	DriveError error
	AuthCode   string
}

func (e *AuthTokenFailed) HumanReadableError(debug bool) error {
	log.Errorf("AuthTokenFailed: %s", e.AuthCode)
	if debug {
		return e.DriveError
	} else {
		return nil
	}
}

// Remote filesystem errors
type QueryFailed struct {
	DriveError error
	Query      string
}

func (e *QueryFailed) HumanReadableError(debug bool) error {
	log.Errorf("QueryFailed: %v", e.DriveError)
	if debug {
		return fmt.Errorf("QueryFailed: %s %w", e.Query, e.DriveError)
	} else {
		return nil
	}
}

type FileNotFound struct {
	DriveError error
	File       string
	Path       string
}

func (e *FileNotFound) HumanReadableError(debug bool) error {
	log.Errorf("FileNotFound: file %s in %s", e.File, e.Path)
	if debug {
		return e.DriveError
	} else {
		return nil
	}
}

type FolderNotFound struct {
	DriveError error
	Path       string
}

func (e *FolderNotFound) HumanReadableError(debug bool) error {
	log.Errorf("FolderNotFound: %s", e.Path)
	if debug {
		return e.DriveError
	} else {
		return nil
	}
}

type DuplicateFile struct {
	File string
	Path string
}

func (e *DuplicateFile) HumanReadableError(debug bool) error {
	log.Errorf("DuplicateFile: file %s already exists in %s", e.File, e.Path)
	if debug {
		return fmt.Errorf("DuplicateFile: file %s already exists in %s", e.File, e.Path)
	} else {
		return nil
	}
}

type OverwriteFailed struct {
	DriveError error
	File       string
	Path       string
}

func (e *OverwriteFailed) HumanReadableError(debug bool) error {
	log.Errorf("OverwriteFailed: file %s in %s", e.File, e.Path)
	if debug {
		return e.DriveError
	} else {
		return nil
	}
}

type CreateFailed struct {
	DriveError error
	File       string
	Path       string
}

func (e *CreateFailed) HumanReadableError(debug bool) error {
	log.Errorf("CreateFailed: file %s in %s", e.File, e.Path)
	if debug {
		return e.DriveError
	} else {
		return nil
	}
}

// Local filesystem errors
type ReadDirFailed struct {
	OSError error
	Path    string
}

func (e *ReadDirFailed) HumanReadableError(debug bool) error {
	log.Errorf("ReadDirFailed: %s", e.Path)
	if debug {
		return e.OSError
	} else {
		return nil
	}
}

type CreateFileFailed struct {
	OSError error
	File    string
}

func (e *CreateFileFailed) HumanReadableError(debug bool) error {
	log.Errorf("CreateFileFailed: %s", e.File)
	if debug {
		return e.OSError
	} else {
		return nil
	}
}

type WriteFileFailed struct {
	OSError error
	File    string
}

func (e *WriteFileFailed) HumanReadableError(debug bool) error {
	log.Errorf("WriteFileFailed: %s", e.File)
	if debug {
		return e.OSError
	} else {
		return nil
	}
}

type CreateDirFailed struct {
	OSError error
	Dir     string
}

func (e *CreateDirFailed) HumanReadableError(debug bool) error {
	log.Errorf("CreateDirFailed: %s", e.Dir)
	if debug {
		return e.OSError
	} else {
		return nil
	}
}

type ReadFileInfoFailed struct {
	OSError error
	File    string
	Path    string
}

func (e *ReadFileInfoFailed) HumanReadableError(debug bool) error {
	log.Errorf("ReadFileInfoFailed: for %s in %s", e.File, e.Path)
	if debug {
		return e.OSError
	} else {
		return nil
	}
}

type OpenFileFailed struct {
	OSError error
	File    string
}

func (e *OpenFileFailed) HumanReadableError(debug bool) error {
	log.Errorf("OpenFileFailed: %s", e.File)
	if debug {
		return e.OSError
	} else {
		return nil
	}
}

type RemoveFileFailed struct {
	OSError error
	File    string
}

func (e *RemoveFileFailed) HumanReadableError(debug bool) error {
	log.Errorf("RemoveFileFailed: %s", e.File)
	if debug {
		return e.OSError
	} else {
		return nil
	}
}

type CopyFileFailed struct {
	OSError error
	File    string
}

func (e *CopyFileFailed) HumanReadableError(debug bool) error {
	log.Errorf("CopyFileFailed: %s", e.File)
	if debug {
		return e.OSError
	} else {
		return nil
	}
}

type ChtimeFailed struct {
	OSError error
	File    string
}

func (e *ChtimeFailed) HumanReadableError(debug bool) error {
	log.Errorf("ChtimeFailed: %s", e.File)
	if debug {
		return e.OSError
	} else {
		return nil
	}
}

type ReadFileFailed struct {
	OSError error
	File    string
}

func (e *ReadFileFailed) HumanReadableError(debug bool) error {
	log.Errorf("ReadFileFailed: %s", e.File)
	if debug {
		return e.OSError
	} else {
		return nil
	}
}

type MimeTypeFailed struct {
	OSError error
	File    string
}

func (e *MimeTypeFailed) HumanReadableError(debug bool) error {
	log.Errorf("MimeTypeFailed: %s", e.File)
	if debug {
		return e.OSError
	} else {
		return nil
	}
}

type Md5Failed struct {
	OSError error
	File    string
}

func (e *Md5Failed) HumanReadableError(debug bool) error {
	log.Errorf("Md5Failed: %s", e.File)
	if debug {
		return e.OSError
	} else {
		return nil
	}
}

type WriteTokenFailed struct {
	OSError error
	File    string
}

func (e *WriteTokenFailed) HumanReadableError(debug bool) error {
	log.Errorf("WriteTokenFailed: %s", e.File)
	if debug {
		return e.OSError
	} else {
		return nil
	}
}

type ParseTokenFailed struct {
	OSError error
}

func (e *ParseTokenFailed) HumanReadableError(debug bool) error {
	log.Errorf("ParseTokenFailed: %v", e.OSError)
	if debug {
		return e.OSError
	} else {
		return nil
	}
}

type TokenDecodeFailed struct {
	OSError error
}

func (e *TokenDecodeFailed) HumanReadableError(debug bool) error {
	log.Errorf("TokenDecodeFailed: %v", e.OSError)
	if debug {
		return e.OSError
	} else {
		return nil
	}
}

type JSONMarshalFailed struct {
	OSError error
	Name    string
}

func (e *JSONMarshalFailed) HumanReadableError(debug bool) error {
	log.Errorf("JSONMarshalFailed: %s", e.Name)
	if debug {
		return e.OSError
	} else {
		return nil
	}
}

type JSONUnmarshalFailed struct {
	OSError error
	File    string
}

func (e *JSONUnmarshalFailed) HumanReadableError(debug bool) error {
	log.Errorf("JSONUnmarshalFailed: %s", e.File)
	if debug {
		return e.OSError
	} else {
		return nil
	}
}

// Internal logic errors
type WrongSyncMode struct {
	Mode SyncMode
}

func (e *WrongSyncMode) HumanReadableError(debug bool) error {
	log.Errorf("WrongSyncMode: %v", e.Mode)
	if debug {
		return fmt.Errorf("WrongSyncMode: %v", e.Mode)
	} else {
		return nil
	}
}
