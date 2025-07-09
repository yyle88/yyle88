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

var mutexRewriteFp sync.Mutex //write file one by one

func TestGenAboutMe(t *testing.T) {
	mutexRewriteFp.Lock()
	defer mutexRewriteFp.Unlock()

	keysBundle, _ := i18n_aboutmekeys.LoadI18nFiles()
	require.NotEmpty(t, keysBundle.LanguageTags())

	valsBundle, _ := i18n_aboutmevals.LoadI18nFiles()
	require.NotEmpty(t, valsBundle.LanguageTags())

	for idx, one := range supportedLanguages {
		caseName := fmt.Sprintf("%d-%s", idx, one.LangCode)

		t.Run(caseName, func(t *testing.T) {
			keysLocalizer := i18n.NewLocalizer(keysBundle, one.LangCode)
			valsLocalizer := i18n.NewLocalizer(valsBundle, one.LangCode)

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

			path, ok := obtainFilename(one.ReadmeFileName)
			if ok {
				return
			}
			t.Log(osmustexist.PATH(path))

			if !slices.Contains([]string{"en", "zh"}, one.LangCode) {
				return //还没写翻译继续保持原状等有空再写
			}

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

func TestGenMarkdownTable(t *testing.T) {
	i18nBundle, messageFiles := i18n_message.LoadI18nFiles()
	require.NotEmpty(t, messageFiles)
	require.NotEmpty(t, i18nBundle.LanguageTags())

	for idx, one := range supportedLanguages {
		caseName := fmt.Sprintf("%d-%s", idx, one.LangCode)

		t.Run(caseName, func(t *testing.T) {
			localizer := i18n.NewLocalizer(i18nBundle, one.LangCode)

			GenMarkdownTable(t, &DocGenParam{
				readmeFileName: one.ReadmeFileName,
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

func GenMarkdownTable(t *testing.T, arg *DocGenParam) {
	mutexRewriteFp.Lock()
	defer mutexRewriteFp.Unlock()

	repos := fetchRepos()
	require.NotEmpty(t, repos)

	ptx := utils.NewPTX()

	cardThemes := utils.GetRepoCardThemes()
	require.NotEmpty(t, cardThemes)

	subRepos, repos := splitRepos(repos, 10)

	ptx.Println(`<div align="left">`)
	ptx.Println()
	for _, repo := range subRepos {
		cardLine := makeCardLine(repo, cardThemes[rand.IntN(len(cardThemes))])

		ptx.Println(cardLine)
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
		const stepLimit = 4
		ptx.Println("|" + repeatString(" "+arg.repoTitle+" |", stepLimit))
		ptx.Println("|" + repeatString(" :--: |", stepLimit))
		for start := 0; start < len(repos); start += stepLimit {
			ptx.Print("|")
			for num := 0; num < stepLimit; num++ {
				if idx := start + num; idx < len(repos) {
					repo := repos[idx]
					ptx.Print(makeCustomHeightBadge(repo.Name, repo.Link, colors[idx%len(colors)], 24), " | ")
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

	path, ok := obtainFilename(arg.readmeFileName)
	if ok {
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
	sIdx := slices.Index(contentLines, param.startLine)
	require.Positive(t, sIdx)
	eIdx := slices.Index(contentLines, param.closeLine)
	require.Positive(t, eIdx)

	require.Less(t, sIdx, eIdx)

	content := strings.Join(contentLines[:sIdx+1], "\n") + "\n" + "\n" +
		param.newString + "\n" +
		strings.Join(contentLines[eIdx:], "\n")
	t.Log(content)

	must.Done(os.WriteFile(param.path, []byte(content), 0666))
	t.Log("success")
}

func obtainFilename(shortName string) (string, bool) {
	path := runpath.PARENT.Join(shortName)
	if !osomitexist.IsFile(path) {
		path = runpath.PARENT.Join("locales", shortName)
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

//func makeBadge(repo *yyle88.Repo, colorString string) string {
//	return fmt.Sprintf("[![%s](https://img.shields.io/badge/%s-%s.svg?style=flat&logoColor=white)](%s)", repo.Name, repo.Name, url.QueryEscape(colorString), repo.Link)
//}

func makeCustomHeightBadge(name string, link string, colorString string, height int) string {
	return fmt.Sprintf(`<a href="%s"><img src="https://img.shields.io/badge/%s-%s.svg?style=flat&logoColor=white" height="%d"></a>`, link, strings.ReplaceAll(name, "-", "+"), url.QueryEscape(colorString), height)
}

func repeatString(s string, n int) string {
	var res string
	for i := 0; i < n; i++ {
		res += s
	}
	return res
}

func makeCardLine(repo *yyle88.Repo, cardTheme string) string {
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
