package yyle88

import "fmt"

type LanguageLink struct {
	Name string
	URL  string
}

func (lang *LanguageLink) String() string {
	return fmt.Sprintf("<a href=\"%s\">%s</a>", lang.URL, lang.Name)
}

func (lang *LanguageLink) Strong() string {
	return fmt.Sprintf("<strong>%s</strong>", lang.Name)
}
