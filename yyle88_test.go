package yyle88_test

import (
	"fmt"
	"math/rand/v2"
	"net/url"
	"os"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/done"
	"github.com/yyle88/must"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/osexistpath/osomitexist"
	"github.com/yyle88/printgo"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
	"github.com/yyle88/yyle88"
	"github.com/yyle88/yyle88/internal/utils"
	"github.com/yyle88/yyle88/locales/i18n_aboutmekeys"
	"github.com/yyle88/yyle88/locales/i18n_aboutmevals"
	"github.com/yyle88/yyle88/locales/i18n_message"
)

const username = "yyle88"

// fileWriteMutex ensures files are written one at a time to prevent conflicts
// fileWriteMutex 确保文件逐个写入以防止冲突
var fileWriteMutex sync.Mutex

// TestGenAboutMe_Preview generates "About Me" section for primary languages (English and Chinese)
// This test is for quick debugging and validation
// TestGenAboutMe_Preview 为主要语言（英文和中文）生成"自我介绍"部分
// 此测试用于快速调试和验证
func TestGenAboutMe_Preview(t *testing.T) {
	primaryLanguages := supportedLanguages[:2] // English and 简体中文
	runGenAboutMe(t, primaryLanguages)
}

// TestGenAboutMe_Publish generates "About Me" section for all remaining languages
// TestGenAboutMe_Publish 为所有剩余语言生成"自我介绍"部分
func TestGenAboutMe_Publish(t *testing.T) {
	publishLanguages := supportedLanguages[2:] // All remaining languages
	runGenAboutMe(t, publishLanguages)
}

// runGenAboutMe generates "About Me" section for specified languages
// runGenAboutMe 为指定语言生成"自我介绍"部分
func runGenAboutMe(t *testing.T, languages []*yyle88.LanguageLink) {
	fileWriteMutex.Lock()
	defer fileWriteMutex.Unlock()

	keysBundle, _ := i18n_aboutmekeys.LoadI18nFiles()
	require.NotEmpty(t, keysBundle.LanguageTags())

	valsBundle, _ := i18n_aboutmevals.LoadI18nFiles()
	require.NotEmpty(t, valsBundle.LanguageTags())

	for idx, lang := range languages {
		testName := fmt.Sprintf("%d-%s", idx, lang.LangCode)

		t.Run(testName, func(t *testing.T) {
			keysLocalizer := i18n.NewLocalizer(keysBundle, lang.LangCode)
			valsLocalizer := i18n.NewLocalizer(valsBundle, lang.LangCode)

			ptx := printgo.NewPTS()
			ptx.Println("##", keysLocalizer.MustLocalize(i18n_aboutmekeys.I18nAboutMe()))
			ptx.Println()
			ptx.Fprintf("- 😄 **%s:** %s\n", keysLocalizer.MustLocalize(i18n_aboutmekeys.I18nName()), valsLocalizer.MustLocalize(i18n_aboutmevals.I18nName()))
			ptx.Fprintf("- 🔭 **%s:** %s\n", keysLocalizer.MustLocalize(i18n_aboutmekeys.I18nBorn()), valsLocalizer.MustLocalize(i18n_aboutmevals.I18nBorn()))
			ptx.Fprintf("- 🌱 **%s:** %s\n", keysLocalizer.MustLocalize(i18n_aboutmekeys.I18nGender()), valsLocalizer.MustLocalize(i18n_aboutmevals.I18nGender()))
			ptx.Fprintf("- 👯 **%s:** %s\n", keysLocalizer.MustLocalize(i18n_aboutmekeys.I18nEducation()), valsLocalizer.MustLocalize(i18n_aboutmevals.I18nEducation()))
			ptx.Fprintf("- 💼 **%s:** %s\n", keysLocalizer.MustLocalize(i18n_aboutmekeys.I18nWorkExperience()), valsLocalizer.MustLocalize(i18n_aboutmevals.I18nWorkExperience()))
			ptx.Fprintf("- 📫 **%s:** %s\n", keysLocalizer.MustLocalize(i18n_aboutmekeys.I18nMainLanguage()), valsLocalizer.MustLocalize(i18n_aboutmevals.I18nMainLanguage()))
			ptx.Fprintf("- 💬 **%s:** %s\n", keysLocalizer.MustLocalize(i18n_aboutmekeys.I18nInterests()), valsLocalizer.MustLocalize(i18n_aboutmevals.I18nInterests()))
			ptx.Fprintf("- 🔗 **%s:** %s\n", keysLocalizer.MustLocalize(i18n_aboutmekeys.I18nGithub()), valsLocalizer.MustLocalize(i18n_aboutmevals.I18nGithub()))
			ptx.Fprintf("- 🌟 **%s:** %s\n", keysLocalizer.MustLocalize(i18n_aboutmekeys.I18nNote()), valsLocalizer.MustLocalize(i18n_aboutmevals.I18nNote()))
			ptx.Fprintf("- ⬆️ **%s:** %s\n", keysLocalizer.MustLocalize(i18n_aboutmekeys.I18nThanks()), valsLocalizer.MustLocalize(i18n_aboutmevals.I18nThanks()))

			t.Log(ptx.String())

			path, skip := resolveReadmePath(lang.ReadmeFileName)
			if skip {
				return
			}
			t.Log(osmustexist.PATH(path))

			replaceBetween(t, &replaceBetweenParam{
				path:      path,
				startLine: "<!-- 这是一个注释，它不会在渲染时显示出来，这是自我介绍的起始位置 -->",
				closeLine: "<!-- 这是一个注释，它不会在渲染时显示出来，这是自我介绍的终止位置 -->",
				newString: ptx.String(),
			})
		})
	}
}

type DocGenParam struct {
	readmeFileName string
	tableTitle     string
	repoTitle      string
}

// runGenMarkdownTable generates markdown table for specified languages
// runGenMarkdownTable 为指定语言生成 markdown 表格
func runGenMarkdownTable(t *testing.T, languages []*yyle88.LanguageLink) {
	i18nBundle, messageFiles := i18n_message.LoadI18nFiles()
	require.NotEmpty(t, messageFiles)
	require.NotEmpty(t, i18nBundle.LanguageTags())

	for idx, lang := range languages {
		testName := fmt.Sprintf("%d-%s", idx, lang.LangCode)

		t.Run(testName, func(t *testing.T) {
			localizer := i18n.NewLocalizer(i18nBundle, lang.LangCode)

			GenMarkdownTable(t, &DocGenParam{
				readmeFileName: lang.ReadmeFileName,
				tableTitle: fmt.Sprintf(
					"| **%s** | **%s** |",
					rese.C1(localizer.Localize(i18n_message.I18nRepoTableTitleName())),
					rese.C1(localizer.Localize(i18n_message.I18nRepoTableTitleDesc())),
				),
				repoTitle: rese.C1(localizer.Localize(i18n_message.I18nRepoTableRepoTitle())),
			})
		})
	}
}

// TestGenMarkdownTable_Preview generates markdown table for primary languages (English and Chinese)
// This test is for quick debugging and validation
// TestGenMarkdownTable_Preview 为主要语言（英文和中文）生成项目表格
// 此测试用于快速调试和验证
func TestGenMarkdownTable_Preview(t *testing.T) {
	primaryLanguages := supportedLanguages[:2] // English and 简体中文
	runGenMarkdownTable(t, primaryLanguages)
}

// TestGenMarkdownTable_Publish generates markdown table for all remaining languages
// TestGenMarkdownTable_Publish 为所有剩余语言生成项目表格
func TestGenMarkdownTable_Publish(t *testing.T) {
	publishLanguages := supportedLanguages[2:] // All remaining languages
	runGenMarkdownTable(t, publishLanguages)
}

func GenMarkdownTable(t *testing.T, arg *DocGenParam) {
	fileWriteMutex.Lock()
	defer fileWriteMutex.Unlock()

	repos := fetchRepos()
	require.NotEmpty(t, repos)

	ptx := utils.NewPTX()

	cardThemes := utils.GetRepoCardThemes()
	require.NotEmpty(t, cardThemes)

	subRepos, repos := splitRepos(repos, 10)

	ptx.Println(`<div align="left">`)
	ptx.Println()
	for _, repo := range subRepos {
		cardMarkdown := generateRepoCard(repo, cardThemes[rand.IntN(len(cardThemes))])

		ptx.Println(cardMarkdown)
		ptx.Println()
	}
	ptx.Println(`</div>`)
	ptx.Println()

	colors := utils.GetBadgeColors()
	require.NotEmpty(t, colors)

	rand.Shuffle(len(colors), func(i, j int) {
		colors[i], colors[j] = colors[j], colors[i]
	})

	if len(repos) > 0 {
		ptx.Println()
		ptx.Println(`<div align="left">`)
		ptx.Println()
		const columnsPerRow = 4
		ptx.Println("|" + strings.Repeat(" "+arg.repoTitle+" |", columnsPerRow))
		ptx.Println("|" + strings.Repeat(" :--: |", columnsPerRow))
		for start := 0; start < len(repos); start += columnsPerRow {
			ptx.Print("|")
			for num := 0; num < columnsPerRow; num++ {
				if idx := start + num; idx < len(repos) {
					repo := repos[idx]
					ptx.Print(generateBadgeWithHeight(repo.Name, repo.Link, colors[idx%len(colors)], 24), " | ")
				} else {
					ptx.Print("-", " | ")
				}
			}
			ptx.Println()
		}
		ptx.Println()
		ptx.Println(`</div>`)
		ptx.Println()
	}

	t.Log(ptx.String())

	path, skip := resolveReadmePath(arg.readmeFileName)
	if skip {
		return
	}
	t.Log(osmustexist.PATH(path))

	replaceBetween(t, &replaceBetweenParam{
		path:      path,
		startLine: "<!-- 这是一个注释，它不会在渲染时显示出来，这是项目列表的起始位置 -->",
		closeLine: "<!-- 这是一个注释，它不会在渲染时显示出来，这是项目列表的终止位置 -->",
		newString: ptx.String(),
	})
}

type replaceBetweenParam struct {
	path      string
	startLine string
	closeLine string
	newString string
}

func replaceBetween(t *testing.T, param *replaceBetweenParam) {
	text := string(done.VAE(os.ReadFile(param.path)).Nice())
	t.Log(text)

	contentLines := strings.Split(text, "\n")
	startIdx := slices.Index(contentLines, param.startLine)
	require.Positive(t, startIdx)
	endIdx := slices.Index(contentLines, param.closeLine)
	require.Positive(t, endIdx)

	require.Less(t, startIdx, endIdx)

	content := strings.Join(contentLines[:startIdx+1], "\n") + "\n" + "\n" +
		param.newString + "\n" +
		strings.Join(contentLines[endIdx:], "\n")
	t.Log(content)

	must.Done(os.WriteFile(param.path, []byte(content), 0666))
	t.Log("success")
}

// resolveReadmePath resolves the full path for a README file
// Returns the path and whether to skip (if file not found)
//
// resolveReadmePath 解析 README 文件的完整路径
// 返回路径以及是否跳过（如果文件未找到）
func resolveReadmePath(filename string) (path string, skip bool) {
	path = runpath.PARENT.Join(filename)
	if !osomitexist.IsFile(path) {
		path = runpath.PARENT.Join("locales", filename)
	}
	if !osomitexist.IsFile(path) {
		return "", true
	}
	return path, false
}

var reposSingleton []*yyle88.Repo
var onceFetchRepos sync.Once

func fetchRepos() []*yyle88.Repo {
	onceFetchRepos.Do(func() {
		reposSingleton = done.VAE(yyle88.GetGithubRepos(username)).Nice()
	})
	return reposSingleton
}

func splitRepos(repos []*yyle88.Repo, subSize int) ([]*yyle88.Repo, []*yyle88.Repo) {
	idx := min(subSize, len(repos))
	return repos[:idx], repos[idx:]
}

// generateBadgeWithHeight creates a custom height badge HTML
// Used for repo links in markdown tables
//
// generateBadgeWithHeight 创建自定义高度的徽章 HTML
// 用于 markdown 表格中的仓库链接
func generateBadgeWithHeight(name, link, colorString string, height int) string {
	return fmt.Sprintf(`<a href="%s"><img src="https://img.shields.io/badge/%s-%s.svg?style=flat&logoColor=white" height="%d"></a>`,
		link, strings.ReplaceAll(name, "-", "+"), url.QueryEscape(colorString), height)
}

// generateRepoCard creates a GitHub repo card markdown with custom theme
// Uses github-readme-stats API for rendering
//
// generateRepoCard 创建带自定义主题的 GitHub 仓库卡片 markdown
// 使用 github-readme-stats API 进行渲染
func generateRepoCard(repo *yyle88.Repo, cardTheme string) string {
	const templateLine = "[![Readme Card](https://github-readme-stats.vercel.app/api/pin/?username={{ username }}&repo={{ repo_name }}&theme={{ card_theme }}&unique={{ unique_uuid }})]({{ repo_link }})"

	rep := strings.NewReplacer(
		"{{ username }}", username,
		"{{ repo_name }}", repo.Name,
		"{{ card_theme }}", cardTheme,
		"{{ unique_uuid }}", uuid.New().String(),
		"{{ repo_link }}", repo.Link,
	)
	cardLine := rep.Replace(templateLine)
	return cardLine
}
