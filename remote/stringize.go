package remote

import (
	"fmt"
	"go-drive/ui"
	"regexp"
	"strings"

	"github.com/dustin/go-humanize"
)

func (r File) Stringize(path, format string) string {
	icon, ok := ui.MimeIcons[r.MimeType]
	if !ok {
		icon = ui.MimeIcons["other"]
	}
	if format == "" {
		format = "%icon% %sync% %shared% %name% %size%"
	}
	if r.MimeType == "application/vnd.google-apps.folder" {
		re := regexp.MustCompile(`\s*[\[\(\{]*%size%[\]\)\}]*\s*`)
		format = re.ReplaceAllString(format, "")
	}
	item_str := format
	item_str = strings.ReplaceAll(item_str, "%icon%", icon)
	item_str = strings.ReplaceAll(item_str, "%size%", fmt.Sprintf("[%s]", humanize.Bytes(uint64(r.Size))))
	item_str = strings.ReplaceAll(item_str, "%name%", r.Name)
	re_shared := regexp.MustCompile(`\s*[\[\(\{]*%shared%[\]\)\}]*\s*`)
	if !r.OwnedByMe {
		item_str = re_shared.ReplaceAllString(item_str, "  ")
	} else {
		item_str = re_shared.ReplaceAllString(item_str, " ")
	}
	if path == "" {
		path = "/"
	}
	re_sync := regexp.MustCompile(`\s*[\[\(\{]*%sync%[\]\)\}]*\s*`)
	item_str = re_sync.ReplaceAllString(item_str, " ")
	// file_mirror := FileMirror{
	// 	RemoteId:  item.Id,
	// 	LocalPath: filepath.Join(path, item.Name),
	// }
	// synced := false
	// if _, ok := syncmap.Mirrors[file_mirror]; ok {
	// 	synced = true
	// }
	// re_sync := regexp.MustCompile(`\s*[\[\(\{]*%sync%[\]\)\}]*\s*`)
	// if synced {
	// 	item_str = re_sync.ReplaceAllString(item_str, " 󰅟 ")
	// } else if item.MimeType == "application/vnd.google-apps.folder" {
	// 	item_str = re_sync.ReplaceAllString(item_str, " ")
	// } else {
	// 	item_str = re_sync.ReplaceAllString(item_str, "  ")
	// }
	// item_str += fmt.Sprintf(" %s", item.ModifiedTime)
	return item_str
}
