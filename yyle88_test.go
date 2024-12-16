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
	//otherDesc string
	closeWith string
}

func TestGenMarkdown(t *testing.T) {
	const username = "yyle88"

	GenMarkdownTable(t, username, &DocGenParam{
		shortName: "README.md",
		startWith: "Here are some of my key projects:",
		titleLine: "| **Repo Name** | **Description** |",
		//otherDesc: "OTHER-PROJECTS:",
		closeWith: "**Explore and star my projects. Your support means a lot!**",
	})
}

func TestGenMarkdownZhHans(t *testing.T) {
	const username = "yyle88"

	GenMarkdownTable(t, username, &DocGenParam{
		shortName: "README.zh.md",
		startWith: "这是我的项目：",
		titleLine: "| 项目名称 | 项目描述 |",
		//otherDesc: "其它项目：",
		closeWith: "给我星星谢谢。",
	})
}

func GenMarkdownTable(t *testing.T, username string, arg *DocGenParam) {
	repos := done.VAE(yyle88.GetGithubRepos(username)).Nice()

	ptx := utils.NewPTX()

	// collect themes from: https://github.com/anuraghazra/github-readme-stats/blob/0810a1d8446902f0267a151302388e5ed8373aa2/themes/README.md?plain=1#L50
	cardThemes := []string{"default_repocard", "transparent", "shadow_red", "shadow_green", "shadow_blue", "dark", "radical", "merko", "gruvbox", "gruvbox_light", "tokyonight", "onedark", "cobalt", "synthwave", "highcontrast", "dracula", "prussian", "monokai", "vue", "vue-dark", "shades-of-purple", "nightowl", "buefy", "blue-green", "algolia", "great-gatsby", "darcula", "bear", "solarized-dark", "solarized-light", "chartreuse-dark", "nord", "gotham", "material-palenight", "graywhite", "vision-friendly-dark", "ayu-mirage", "midnight-purple", "calm", "flag-india", "omni", "react", "jolly", "maroongold", "yeblu", "blueberry", "slateorange", "kacho_ga", "outrun", "ocean_dark", "city_lights", "github_dark", "github_dark_dimmed", "discord_old_blurple", "aura_dark", "panda", "noctis_minimus", "cobalt2", "swift", "aura", "apprentice", "moltack", "codeSTACKr", "rose_pine", "catppuccin_latte", "catppuccin_mocha", "date_night", "one_dark_pro", "rose", "holi", "neon", "blue_navy", "calm_pink", "ambient_gradient"}

	subRepos, repos := splitRepos(repos, 5)
	for _, repo := range subRepos {
		const templateLine = "[![Readme Card](https://github-readme-stats.vercel.app/api/pin/?username={{ username }}&repo={{ repo_name }}&theme={{ card_theme }})]({{ repo_link }})"

		rep := strings.NewReplacer(
			"{{ username }}", username,
			"{{ repo_name }}", repo.Name,
			"{{ card_theme }}", cardThemes[rand.IntN(len(cardThemes))],
			"{{ repo_link }}", repo.Link,
		)

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
		//ptx.Println(arg.otherDesc)
		//ptx.Println()
		const stepLimit = 4
		ptx.Println("|" + repeatString(" |", stepLimit))
		ptx.Println("|" + repeatString(" :--: |", stepLimit))
		for start := 0; start < len(repos); start += stepLimit {
			ptx.Print("|")
			for num := 0; num < stepLimit; num++ {
				if idx := start + num; idx < len(repos) {
					ptx.Print(makeBadge(repos[idx], colors[idx%len(colors)]), " | ")
				} else {
					ptx.Print("-", " | ")
				}
			}
			ptx.Println()
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

func repeatString(s string, n int) string {
	var res string
	for i := 0; i < n; i++ {
		res += s
	}
	return res
}
