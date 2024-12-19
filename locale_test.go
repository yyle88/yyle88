package yyle88_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexistpath/osomitexist"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
	"github.com/yyle88/yyle88"
)

var supportedLanguages = []*yyle88.LanguageLink{
	{Name: "English", URL: "./README.md"},       // 英语
	{Name: "简体中文", URL: "./README.zh.md"},       // 简体中文
	{Name: "繁體中文", URL: "./README.zh-Hant.md"},  // 繁体中文
	{Name: "日本語", URL: "./README.ja.md"},        // 日语
	{Name: "Русский", URL: "./README.ru.md"},    // 俄语
	{Name: "Deutsch", URL: "./README.de.md"},    // 德语
	{Name: "Français", URL: "./README.fr.md"},   // 法语
	{Name: "Español", URL: "./README.es.md"},    // 西班牙语
	{Name: "Português", URL: "./README.pt.md"},  // 葡萄牙语
	{Name: "ខ្មែរ", URL: "./README.kh.md"},      // 高棉语
	{Name: "Tiếng Việt", URL: "./README.vi.md"}, // 越南语
}

func TestGenLanguageLinkMarkdown(t *testing.T) {
	mutexRewriteFp.Lock()
	defer mutexRewriteFp.Unlock()

	root := runpath.PARENT.Path()

	var matchedLanguages []*yyle88.LanguageLink
	for _, lang := range supportedLanguages {
		if osomitexist.IsFile(filepath.Join(root, lang.URL)) {
			t.Log(neatjsons.S(lang))
			matchedLanguages = append(matchedLanguages, lang)
		}
	}

	for _, lang := range matchedLanguages {
		path := filepath.Join(root, lang.URL)
		t.Log(path)

		text := string(rese.V1(os.ReadFile(path)))
		t.Log(text)

		var radioLinks []string
		for _, lang2 := range matchedLanguages {
			if lang2.URL == lang.URL {
				radioLinks = append(radioLinks, lang2.Strong())
			} else {
				radioLinks = append(radioLinks, lang2.String())
			}
		}

		contentLines := strings.Split(text, "\n")
		sIdx := slices.Index(contentLines, "<!-- 这是一个注释，它不会在渲染时显示出来，这是语言选择的起始位置 -->")
		require.Positive(t, sIdx)
		eIdx := slices.Index(contentLines, "<!-- 这是一个注释，它不会在渲染时显示出来，这是语言选择的终止位置 -->")
		require.Positive(t, eIdx)

		require.Less(t, sIdx, eIdx)

		content := strings.Join(contentLines[:sIdx+1], "\n") + "\n" + "\n" +
			`<h4 align="center">` + strings.Join(radioLinks, " | ") + `</h4>` + "\n" + "\n" +
			strings.Join(contentLines[eIdx:], "\n")
		t.Log(content)

		must.Done(os.WriteFile(path, []byte(content), 0666))
		t.Log("success")
	}
}
