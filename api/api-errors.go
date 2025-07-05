package api

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
	drive_error error
}

func (e *GoogleDriveError) HumanReadableError(debug bool) error {
	log.Error("GoogleDriveError: %v", e.drive_error)
	if debug {
		return e.drive_error
	} else {
		return nil
	}
}

type ParseCredentialsFailed struct {
	drive_error error
}

func (e *ParseCredentialsFailed) HumanReadableError(debug bool) error {
	log.Error("ParseCredentialsFailed: %v", e.drive_error)
	if debug {
		return e.drive_error
	} else {
		return nil
	}
}

type AuthTokenFailed struct {
	drive_error error
	authcode    string
}

func (e *AuthTokenFailed) HumanReadableError(debug bool) error {
	log.Error("AuthTokenFailed: %s", e.authcode)
	if debug {
		return e.drive_error
	} else {
		return nil
	}
}

// Remote filesystem errors
type QueryFailed struct {
	drive_error error
	query       string
}

func (e *QueryFailed) HumanReadableError(debug bool) error {
	log.Error("QueryFailed: %v", e.drive_error)
	if debug {
		return e.drive_error
	} else {
		return nil
	}
}

type FileNotFound struct {
	drive_error error
	file        string
	path        string
}

func (e *FileNotFound) HumanReadableError(debug bool) error {
	log.Error("FileNotFound: file %s in %s", e.file, e.path)
	if debug {
		return e.drive_error
	} else {
		return nil
	}
}

type FolderNotFound struct {
	drive_error error
	path        string
}

func (e *FolderNotFound) HumanReadableError(debug bool) error {
	log.Error("FolderNotFound: %s", e.path)
	if debug {
		return e.drive_error
	} else {
		return nil
	}
}

type DuplicateFile struct {
	file string
	path string
}

func (e *DuplicateFile) HumanReadableError(debug bool) error {
	log.Error("DuplicateFile: file %s already exists in %s", e.file, e.path)
	if debug {
		return fmt.Errorf("DuplicateFile: file %s already exists in %s", e.file, e.path)
	} else {
		return nil
	}
}

type OverwriteFailed struct {
	drive_error error
	file        string
	path        string
}

func (e *OverwriteFailed) HumanReadableError(debug bool) error {
	log.Error("OverwriteFailed: file %s in %s", e.file, e.path)
	if debug {
		return e.drive_error
	} else {
		return nil
	}
}

type CreateFailed struct {
	drive_error error
	file        string
	path        string
}

func (e *CreateFailed) HumanReadableError(debug bool) error {
	log.Error("CreateFailed: file %s in %s", e.file, e.path)
	if debug {
		return e.drive_error
	} else {
		return nil
	}
}

// Local filesystem errors
type OpenFileFailed struct {
	os_error error
	file     string
}

func (e *OpenFileFailed) HumanReadableError(debug bool) error {
	log.Error("OpenFileFailed: %s", e.file)
	if debug {
		return e.os_error
	} else {
		return nil
	}
}

type ReadFileFailed struct {
	os_error error
	file     string
}

func (e *ReadFileFailed) HumanReadableError(debug bool) error {
	log.Error("ReadFileFailed: %s", e.file)
	if debug {
		return e.os_error
	} else {
		return nil
	}
}

type WriteTokenFailed struct {
	os_error error
	file     string
}

func (e *WriteTokenFailed) HumanReadableError(debug bool) error {
	log.Error("WriteTokenFailed: %s", e.file)
	if debug {
		return e.os_error
	} else {
		return nil
	}
}

type ParseTokenFailed struct {
	os_error error
}

func (e *ParseTokenFailed) HumanReadableError(debug bool) error {
	log.Error("ParseTokenFailed: %v", e.os_error)
	if debug {
		return e.os_error
	} else {
		return nil
	}
}

type TokenDecodeFailed struct {
	os_error error
}

func (e *TokenDecodeFailed) HumanReadableError(debug bool) error {
	log.Error("TokenDecodeFailed: %v", e.os_error)
	if debug {
		return e.os_error
	} else {
		return nil
	}
}

// Internal logic errors
type WrongUploadMode struct {
	mode UploadMode
}

func (e *WrongUploadMode) HumanReadableError(debug bool) error {
	log.Error("WrongUploadMode: %v", e.mode)
	if debug {
		return fmt.Errorf("WrongUploadMode: %v", e.mode)
	} else {
		return nil
	}
}
