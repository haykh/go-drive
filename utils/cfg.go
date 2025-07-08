package utils

type SyncMode int

const (
	RaiseIfDuplicate SyncMode = iota
	SkipDuplicates
	Overwrite
)
