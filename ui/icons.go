package ui

import "github.com/charmbracelet/lipgloss"

var MimeIcons = map[string]string{
	"directory":                                "",
	"application/vnd.google-apps.folder":       "",
	"application/vnd.google.colaboratory":      "",
	"application/vnd.google-apps.document":     "󰈙",
	"application/vnd.google-apps.spreadsheet":  "󰧷",
	"application/vnd.google-apps.presentation": "󰈩",
	"application/pdf":                          "",
	"application/msword":                       "",
	"application/zip":                          "",
	"application/mathematica":                  "󰿈",
	"application/x-gzip":                       "",
	"application/x-xcf":                        "",
	"application/octet-stream":                 "",
	"text/plain; charset=utf-8":                "",
	"text/x-python":                            "",
	"application/json":                         "",
	"application/x-sharedlib":                  "",
	"application/x-object":                     "",
	"video/mp4":                                "",
	"image/png":                                "",
	"image/jpeg":                               "",
	"image/gif":                                "󰈟",
	"image/x-photoshop":                        "",
	"application/illustrator":                  "",
	"text/csv":                                 "",
	"text/xml":                                 "󰗀",
	"other":                                    "",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": "",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   "",
}

var StatusIcons = map[string]string{
	"synced":  "󰅟",
	"remote":  "",
	"local":   "",
	"syncing": "󰘿",
	"shared":  "",
}

var (
	SyncedStyle  = WithForeground("2")
	RemoteStyle  = WithForeground("3")
	LocalStyle   = WithForeground("4")
	SyncingStyle = WithForeground("1")

	MimeIconStyle         = WithForeground("8")
	MimeIconSelectedStyle = WithForeground("170")

	FileNameStyle         = lipgloss.NewStyle()
	FileNameSelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Underline(true)

	FileSizeStyle = WithForeground("8")
)
