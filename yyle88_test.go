package yyle88_test

import (
	"fmt"
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
	closeWith string
}

func TestGenMarkdown(t *testing.T) {
	const username = "yyle88"

	const shortName = "README.md"
	const startWith = "Here are some of my key projects:"
	const titleLine = "| **Project Name** | **Description** |"
	const closeWith = "**Explore and star my projects. Your support means a lot!**"

	GenMarkdownTable(t, username, &DocGenParam{
		shortName: shortName,
		startWith: startWith,
		titleLine: titleLine,
		closeWith: closeWith,
	})
}

func TestGenMarkdownZhHans(t *testing.T) {
	const username = "yyle88"

	const shortName = "README.zh.md"
	const startWith = "这是我的项目："
	const titleLine = "| 项目名称 | 项目描述 |"
	const closeWith = "给我星星谢谢。"

	GenMarkdownTable(t, username, &DocGenParam{
		shortName: shortName,
		startWith: startWith,
		titleLine: titleLine,
		closeWith: closeWith,
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

	subRepos, repos = splitRepos(repos, 5)
	if len(subRepos) > 0 {
		ptx.Println(arg.titleLine)
		ptx.Println("|-------------------------------------------------|--------|")
		for _, repo := range subRepos {
			ptx.Println(fmt.Sprintf("| [%s](%s) | %s |", repo.Name, repo.Link, strings.ReplaceAll(repo.Desc, "|", "-")))
		}
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
