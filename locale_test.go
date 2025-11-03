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
	{LangName: "English", ReadmeFileName: "./README.md", LangCode: "en"},             // 英语
	{LangName: "简体中文", ReadmeFileName: "./README.zh.md", LangCode: "zh"},             // 简体中文
	{LangName: "繁體中文", ReadmeFileName: "./README.zh-Hant.md", LangCode: "zh-Hant"},   // 繁体中文
	{LangName: "日本語", ReadmeFileName: "./README.ja.md", LangCode: "ja"},              // 日语
	{LangName: "Русский", ReadmeFileName: "./README.ru.md", LangCode: "ru"},          // 俄语
	{LangName: "Deutsch", ReadmeFileName: "./README.de.md", LangCode: "de"},          // 德语
	{LangName: "Français", ReadmeFileName: "./README.fr.md", LangCode: "fr"},         // 法语
	{LangName: "Español", ReadmeFileName: "./README.es.md", LangCode: "es"},          // 西班牙语
	{LangName: "Português", ReadmeFileName: "./README.pt.md", LangCode: "pt"},        // 葡萄牙语
	{LangName: "Tiếng Việt", ReadmeFileName: "./README.vi.md", LangCode: "vi"},       // 越南语
	{LangName: "ខ្មែរ", ReadmeFileName: "./README.kh.md", LangCode: "kh"},            // 高棉语
	{LangName: "한국어", ReadmeFileName: "./README.ko.md", LangCode: "ko"},              // 韩语
	{LangName: "Türkçe", ReadmeFileName: "./README.tr.md", LangCode: "tr"},           // 土耳其语
	{LangName: "Polski", ReadmeFileName: "./README.pl.md", LangCode: "pl"},           // 波兰语
	{LangName: "Italiano", ReadmeFileName: "./README.it.md", LangCode: "it"},         // 意大利语
	{LangName: "العربية", ReadmeFileName: "./README.ar.md", LangCode: "ar"},          // 阿拉伯语
	{LangName: "فارسی", ReadmeFileName: "./README.fa.md", LangCode: "fa"},            // 波斯语
	{LangName: "Čeština", ReadmeFileName: "./README.cs.md", LangCode: "cs"},          // 捷克语
	{LangName: "Українська", ReadmeFileName: "./README.uk.md", LangCode: "uk"},       // 乌克兰语
	{LangName: "Nederlands", ReadmeFileName: "./README.nl.md", LangCode: "nl"},       // 荷兰语
	{LangName: "हिन्दी", ReadmeFileName: "./README.hi.md", LangCode: "hi"},           // 印地语
	{LangName: "ภาษาไทย", ReadmeFileName: "./README.th.md", LangCode: "th"},          // 泰语
	{LangName: "Bahasa Indonesia", ReadmeFileName: "./README.id.md", LangCode: "id"}, // 印尼语
	{LangName: "Bahasa Melayu", ReadmeFileName: "./README.ms.md", LangCode: "ms"},    // 马来语
	{LangName: "Filipino", ReadmeFileName: "./README.ph.md", LangCode: "ph"},         // 菲律宾语
	{LangName: "বাংলা", ReadmeFileName: "./README.bn.md", LangCode: "bn"},            // 孟加拉语
}

func TestMoveReadmeIntoLocales(t *testing.T) {
	fileWriteMutex.Lock()
	defer fileWriteMutex.Unlock()

	root := runpath.PARENT.Path()
	for idx, lang := range supportedLanguages {
		readmePath := filepath.Join(root, lang.ReadmeFileName)
		localePath := filepath.Join(root, "locales", lang.ReadmeFileName)
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
	fileWriteMutex.Lock()
	defer fileWriteMutex.Unlock()

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
			filepath.Join(root, lang.ReadmeFileName),
			filepath.Join(root, "locales", lang.ReadmeFileName),
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
	fileWriteMutex.Lock()
	defer fileWriteMutex.Unlock()

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
				radioLinks = append(radioLinks, next.LangLink.StrongLangName())
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
