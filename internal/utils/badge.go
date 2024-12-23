package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func GetBadgeColors() []string {
	return []string{
		"#FF5733",
		"#91C4A4",
		"#7D4B91",
		"#35A8D5",
		"#F2D330",
		"#F09F3B",
		"#F7931E",
		"#95C59D",
		"#7D5E7F",
		"#8A2BE2",
		"#FF6347",
		"#FF1493",
		"#32CD32",
		"#20B2AA",
		"#FFD700",
		"#DC143C",
		"#FF4500",
		"#2E8B57",
		"#3CB371",
		"#ADFF2F",
	}
}

func MakeCustomSizeBadge(name string, link string, colorString string, height int, width int) string {
	msg := fmt.Sprintf(`<a href="%s"><img src="https://img.shields.io/badge/%s-%s.svg?style=flat&logoColor=white"`, link, strings.ReplaceAll(name, "-", "+"), url.QueryEscape(colorString))
	if height > 0 {
		msg += fmt.Sprintf(` height="%d"`, height)
	}
	if width > 0 {
		msg += fmt.Sprintf(` width="%d"`, width)
	}
	msg += `></a>`
	return msg
}
