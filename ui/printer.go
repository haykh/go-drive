package ui

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/dustin/go-humanize"
	"google.golang.org/api/drive/v3"
)

var Icons = map[string]string{
	"application/vnd.google-apps.folder":       "󰉋",
	"application/vnd.google.colaboratory":      "",
	"application/vnd.google-apps.document":     "󰈙",
	"application/vnd.google-apps.spreadsheet":  "󰧷",
	"application/vnd.google-apps.presentation": "󰈩",
	"application/pdf":                          "",
	"application/msword":                       "",
	"image/png":                                "",
	"image/jpeg":                               "",
	"image/gif":                                "󰈟",
	"other":                                    "",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": "",
}

func StringizeItem(item *drive.File, format string) string {
	icon, ok := Icons[item.MimeType]
	if !ok {
		icon = Icons["other"]
	}
	if format == "" {
		format = "%icon% %name% %size%"
	}
	if item.MimeType == "application/vnd.google-apps.folder" {
		re := regexp.MustCompile(`\s*[\[\(\{]*%size%[\]\)\}]*\s*`)
		format = re.ReplaceAllString(format, "")
	}
	item_str := format
	item_str = strings.ReplaceAll(item_str, "%icon%", icon)
	item_str = strings.ReplaceAll(item_str, "%size%", fmt.Sprintf("[%s]", humanize.Bytes(uint64(item.Size))))
	item_str = strings.ReplaceAll(item_str, "%name%", item.Name)
	re := regexp.MustCompile(`\s*[\[\(\{]*%shared%[\]\)\}]*\s*`)
	if !item.OwnedByMe {
		item_str = re.ReplaceAllString(item_str, "  ")
	} else {
		item_str = re.ReplaceAllString(item_str, " ")
	}
	return item_str
}

func SortFileList(items *drive.FileList) *drive.FileList {
	sorted_items := items.Files
	sort.Slice(sorted_items, func(i, j int) bool {
		if items.Files[i].MimeType == items.Files[j].MimeType {
			return items.Files[i].Name < items.Files[j].Name
		} else if items.Files[i].MimeType == "application/vnd.google-apps.folder" {
			return true
		} else if items.Files[j].MimeType == "application/vnd.google-apps.folder" {
			return false
		} else if items.Files[i].MimeType == "application/pdf" {
			return true
		} else if items.Files[j].MimeType == "application/pdf" {
			return false
		}
		if _, ok := Icons[items.Files[i].MimeType]; !ok {
			return false
		}
		if _, ok := Icons[items.Files[j].MimeType]; !ok {
			return true
		}
		return items.Files[i].MimeType < items.Files[j].MimeType
	})
	items.Files = sorted_items
	return items
}

func StringizeItemList(items *drive.FileList, format string) []string {
	sorted_items := SortFileList(items).Files
	items_str := []string{}
	for item := range sorted_items {
		items_str = append(items_str, StringizeItem(sorted_items[item], format))
	}
	return items_str
}
