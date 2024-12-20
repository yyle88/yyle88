package yyle88

import (
	"fmt"
	"path/filepath"
)

type LanguageLink struct {
	Name string
	URL  string
}

func (lang *LanguageLink) String() string {
	return CreateLink(lang.URL, lang.Name)
}

func (lang *LanguageLink) Strong() string {
	return fmt.Sprintf("<strong>%s</strong>", lang.Name)
}

type LangLinkPath struct {
	LangLink *LanguageLink
	Path     string
}

func (a *LangLinkPath) CreateLink(parentPath string) string {
	return CreateLink(filepath.Join(parentPath, a.LangLink.URL), a.LangLink.Name)
}

func CreateLink(link string, name string) string {
	return fmt.Sprintf("<a href=\"%s\">%s</a>", link, name)
}
