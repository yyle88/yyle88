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

func TestGenMarkdown(t *testing.T) {
	username := "yyle88"
	repos := done.VAE(yyle88.GetGithubRepos(username)).Nice()

	repos = repos[:min(5, len(repos))]

	ptx := utils.NewPTX()
	ptx.Println("| 项目名称 | 项目描述 |")
	ptx.Println("|-------------------------------------------------|--------|")
	for _, repo := range repos {
		ptx.Println(fmt.Sprintf("| [%s](%s) | %s |", repo.Name, repo.Link, strings.ReplaceAll(repo.Desc, "|", "-")))
	}

	stb := ptx.String()
	t.Log(stb)

	path := osmustexist.PATH(runpath.PARENT.Join("README.md"))
	t.Log(path)

	text := string(done.VAE(os.ReadFile(path)).Nice())
	t.Log(text)

	sLns := strings.Split(text, "\n")
	sIdx := slices.Index(sLns, "这是我的项目：")
	require.Positive(t, sIdx)
	eIdx := slices.Index(sLns, "给我星星谢谢。")
	require.Positive(t, eIdx)

	require.Less(t, sIdx, eIdx)

	content := strings.Join(sLns[:sIdx+1], "\n") + "\n" + "\n" +
		stb + "\n" +
		strings.Join(sLns[eIdx:], "\n")
	t.Log(content)

	must.Done(os.WriteFile(path, []byte(content), 0666))
	t.Log("success")
}
