package yyle88_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/done"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/osexistpath/osomitexist"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
	"github.com/yyle88/tern"
	"github.com/yyle88/yyle88"
	"github.com/yyle88/yyle88/internal/utils"
)

var supportedLanguages = []*yyle88.LanguageLink{
	{Name: "English", URL: "./README.md"},             // 英语
	{Name: "简体中文", URL: "./README.zh.md"},             // 简体中文
	{Name: "繁體中文", URL: "./README.zh-Hant.md"},        // 繁体中文
	{Name: "日本語", URL: "./README.ja.md"},              // 日语
	{Name: "Русский", URL: "./README.ru.md"},          // 俄语
	{Name: "Deutsch", URL: "./README.de.md"},          // 德语
	{Name: "Français", URL: "./README.fr.md"},         // 法语
	{Name: "Español", URL: "./README.es.md"},          // 西班牙语
	{Name: "Português", URL: "./README.pt.md"},        // 葡萄牙语
	{Name: "Tiếng Việt", URL: "./README.vi.md"},       // 越南语
	{Name: "ខ្មែរ", URL: "./README.kh.md"},            // 高棉语
	{Name: "한국어", URL: "./README.ko.md"},              // 韩国语
	{Name: "Türkçe", URL: "./README.tr.md"},           // 土耳其语
	{Name: "Polski", URL: "./README.pl.md"},           // 波兰语
	{Name: "Italiano", URL: "./README.it.md"},         // 意大利语
	{Name: "العربية", URL: "./README.ar.md"},          // 阿拉伯语
	{Name: "فارسی", URL: "./README.fa.md"},            // 波斯语
	{Name: "Čeština", URL: "./README.cs.md"},          // 捷克语
	{Name: "Українська", URL: "./README.uk.md"},       // 乌克兰语
	{Name: "Nederlands", URL: "./README.nl.md"},       // 荷兰语
	{Name: "हिन्दी", URL: "./README.hi.md"},           // 印地语
	{Name: "ภาษาไทย", URL: "./README.th.md"},          // 泰语
	{Name: "Bahasa Indonesia", URL: "./README.id.md"}, // 印尼语
	{Name: "Bahasa Melayu", URL: "./README.ms.md"},    // 马来语
	{Name: "Filipino", URL: "./README.ph.md"},         // 菲律宾语
	{Name: "বাংলা", URL: "./README.bn.md"},            // 孟加拉语
}

func TestMoveReadmeIntoLocales(t *testing.T) {
	mutexRewriteFp.Lock()
	defer mutexRewriteFp.Unlock()

	root := runpath.PARENT.Path()
	for idx, lang := range supportedLanguages {
		readmePath := filepath.Join(root, lang.URL)
		localePath := filepath.Join(root, "locales", lang.URL)
		if idx < 2 {
			if osomitexist.IsFile(localePath) {
				done.VAE(osexec.Exec("mv", localePath, readmePath)).Done()
			}
		} else {
			if osomitexist.IsFile(readmePath) {
				done.VAE(osexec.Exec("mv", readmePath, localePath)).Done()
			}
		}
	}
}

func TestWriteLocaleMenu(t *testing.T) {
	mutexRewriteFp.Lock()
	defer mutexRewriteFp.Unlock()

	const menuShortName = "LOCALE-MENU.md"
	menuPath := osmustexist.PATH(runpath.PARENT.Join(menuShortName))
	t.Log(menuPath)

	matchedLanguages := caseGetMatchedLanguages(t)

	ptx := utils.NewPTX()
	ptx.Println()
	ptx.Println("<div style=\"text-align: center;\">")
	ptx.Println("<table style=\"margin: 0 auto; text-align: center;\">")
	ptx.Println("<tr><th><strong>LANGUAGE</strong></th></tr>")
	for _, next := range matchedLanguages {
		relativePath := rese.V1(filepath.Rel(filepath.Dir(menuPath), filepath.Dir(next.Path)))
		newLinkString := next.CreateLink(filepath.Join(".", relativePath))
		ptx.Println("<tr><td>" + newLinkString + "</td></tr>")
	}
	ptx.Println("</table>")
	ptx.Println("</div>")
	stb := ptx.String()
	t.Log(stb)

	rewriteLanguageTable(t, menuPath, stb)
}

func caseGetMatchedLanguages(t *testing.T) []*yyle88.LangLinkPath {
	root := runpath.PARENT.Path()

	var matchedLanguages []*yyle88.LangLinkPath
	for _, lang := range supportedLanguages {
		for _, path := range []string{
			filepath.Join(root, lang.URL),
			filepath.Join(root, "locales", lang.URL),
		} {
			if osomitexist.IsFile(path) {
				t.Log(neatjsons.S(lang))
				matchedLanguages = append(matchedLanguages, &yyle88.LangLinkPath{
					LangLink: lang,
					Path:     path,
				})
				break
			}
		}
	}
	return matchedLanguages
}

func TestGenLanguageLinkMarkdown(t *testing.T) {
	mutexRewriteFp.Lock()
	defer mutexRewriteFp.Unlock()

	menuRoot := runpath.PARENT.Path()

	const menuShortName = "LOCALE-MENU.md"
	menuPath := osmustexist.PATH(filepath.Join(menuRoot, menuShortName))
	t.Log(menuPath)

	matchedLanguages := caseGetMatchedLanguages(t)

	const maxOutLangCount = 10

	for _, lang := range matchedLanguages {
		var radioLinks []string

		var meetSamePath = false
		for _, next := range matchedLanguages {
			if next.Path == lang.Path {
				radioLinks = append(radioLinks, next.LangLink.Strong())
				meetSamePath = true
			} else {
				if len(radioLinks)+tern.BVV(meetSamePath, 0, 1) < maxOutLangCount {
					relativePath := rese.V1(filepath.Rel(filepath.Dir(lang.Path), filepath.Dir(next.Path)))
					newLinkString := next.CreateLink(filepath.Join(".", relativePath))
					radioLinks = append(radioLinks, newLinkString)
				}
			}
		}

		relativePath := rese.V1(filepath.Rel(filepath.Dir(lang.Path), menuRoot))

		radioLinks = append(radioLinks, yyle88.CreateLink(filepath.Join(relativePath, menuShortName), "<b>...</b>"))

		stb := `<h4 align="center" style="font-size: 2.0em;">` + strings.Join(radioLinks, " | ") + `</h4>`

		rewriteLanguageTable(t, lang.Path, stb)
	}
}

func rewriteLanguageTable(t *testing.T, path string, stb string) {
	t.Log(path)

	text := string(rese.V1(os.ReadFile(path)))
	t.Log(text)

	contentLines := strings.Split(text, "\n")
	sIdx := slices.Index(contentLines, "<!-- 这是一个注释，它不会在渲染时显示出来，这是语言选择的起始位置 -->")
	require.Positive(t, sIdx)
	eIdx := slices.Index(contentLines, "<!-- 这是一个注释，它不会在渲染时显示出来，这是语言选择的终止位置 -->")
	require.Positive(t, eIdx)

	require.Less(t, sIdx, eIdx)

	content := strings.Join(contentLines[:sIdx+1], "\n") + "\n" + "\n" +
		stb + "\n" + "\n" +
		strings.Join(contentLines[eIdx:], "\n")
	t.Log(content)

	must.Done(os.WriteFile(path, []byte(content), 0666))
	t.Log("success")
}
