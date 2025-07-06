package utils

type UploadMode int

const (
	RaiseIfDuplicate UploadMode = iota
	SkipDuplicates
	Overwrite
)
