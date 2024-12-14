package yyle88_test

import (
	"fmt"
	"math/rand/v2"
	"net/url"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/done"
	"github.com/yyle88/must"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/runpath"
	"github.com/yyle88/yyle88"
	"github.com/yyle88/yyle88/internal/utils"
)

type DocGenParam struct {
	shortName string
	startWith string
	titleLine string
	otherDesc string
	closeWith string
}

func TestGenMarkdown(t *testing.T) {
	const username = "yyle88"

	GenMarkdownTable(t, username, &DocGenParam{
		shortName: "README.md",
		startWith: "Here are some of my key projects:",
		titleLine: "| **Repo Name** | **Description** |",
		otherDesc: "OTHER-PROJECTS:",
		closeWith: "**Explore and star my projects. Your support means a lot!**",
	})
}

func TestGenMarkdownZhHans(t *testing.T) {
	const username = "yyle88"

	GenMarkdownTable(t, username, &DocGenParam{
		shortName: "README.zh.md",
		startWith: "这是我的项目：",
		titleLine: "| 项目名称 | 项目描述 |",
		otherDesc: "其它项目：",
		closeWith: "给我星星谢谢。",
	})
}

func GenMarkdownTable(t *testing.T, username string, arg *DocGenParam) {
	repos := done.VAE(yyle88.GetGithubRepos(username)).Nice()

	ptx := utils.NewPTX()

	subRepos, repos := splitRepos(repos, 5)
	for _, repo := range subRepos {
		const templateLine = "[![Readme Card](https://github-readme-stats.vercel.app/api/pin/?username={{ username }}&repo={{ repo_name }}&theme=algolia)]({{ repo_link }})"

		rep := strings.NewReplacer("{{ username }}", username, "{{ repo_name }}", repo.Name, "{{ repo_link }}", repo.Link)

		ptx.Println(rep.Replace(templateLine))
		ptx.Println()
	}

	colors := []string{"#FF5733", "#91C4A4", "#7D4B91", "#35A8D5", "#F2D330", "#F09F3B", "#F7931E", "#95C59D", "#7D5E7F", "#8A2BE2", "#FF6347", "#FF1493", "#32CD32", "#20B2AA", "#FFD700", "#DC143C", "#FF4500", "#2E8B57", "#3CB371", "#ADFF2F"}
	rand.Shuffle(len(colors), func(i, j int) {
		colors[i], colors[j] = colors[j], colors[i]
	})

	subRepos, repos = splitRepos(repos, 5)
	if len(subRepos) > 0 {
		ptx.Println()
		ptx.Println(arg.titleLine)
		ptx.Println("|-------------------------------------------------|--------|")
		for _, repo := range subRepos {
			ptx.Println(fmt.Sprintf("| %s | %s |", makeBadge(repo, colors[rand.IntN(len(colors))]), strings.ReplaceAll(repo.Desc, "|", "-")))
		}
		ptx.Println()
	}

	if len(repos) > 0 {
		ptx.Println()
		ptx.Println(arg.otherDesc)
		for idx, repo := range repos {
			ptx.Println(makeBadge(repo, colors[idx%len(colors)]))
		}
		ptx.Println()
	}

	stb := ptx.String()
	t.Log(stb)

	path := osmustexist.PATH(runpath.PARENT.Join(arg.shortName))
	t.Log(path)

	text := string(done.VAE(os.ReadFile(path)).Nice())
	t.Log(text)

	sLns := strings.Split(text, "\n")
	sIdx := slices.Index(sLns, arg.startWith)
	require.Positive(t, sIdx)
	eIdx := slices.Index(sLns, arg.closeWith)
	require.Positive(t, eIdx)

	require.Less(t, sIdx, eIdx)

	content := strings.Join(sLns[:sIdx+1], "\n") + "\n" + "\n" +
		stb + "\n" +
		strings.Join(sLns[eIdx:], "\n")
	t.Log(content)

	must.Done(os.WriteFile(path, []byte(content), 0666))
	t.Log("success")
}

func splitRepos(repos []*yyle88.Repo, subSize int) ([]*yyle88.Repo, []*yyle88.Repo) {
	idx := min(subSize, len(repos))
	return repos[:idx], repos[idx:]
}

func makeBadge(repo *yyle88.Repo, colorString string) string {
	return fmt.Sprintf("[![%s](https://img.shields.io/badge/%s-%s.svg?style=flat&logoColor=white)](%s)", repo.Name, repo.Name, url.QueryEscape(colorString), repo.Link)
}
