package yyle88

import (
	"fmt"
	"path/filepath"
)

type LanguageLink struct {
	LangName       string
	ReadmeFileName string
	LangCode       string
}

func (lang *LanguageLink) StrongLangName() string {
	return fmt.Sprintf("<strong>%s</strong>", lang.LangName)
}

type LangLinkPath struct {
	LangLink *LanguageLink
	Path     string
}

func (a *LangLinkPath) CreateLink(parentPath string) string {
	return CreateLink(filepath.Join(parentPath, a.LangLink.ReadmeFileName), a.LangLink.LangName)
}

func CreateLink(link string, name string) string {
	return fmt.Sprintf("<a href=\"%s\">%s</a>", link, name)
}
