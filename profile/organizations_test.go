package profile

import (
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/done"
	"github.com/yyle88/must"
	"github.com/yyle88/mutexmap"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
	"github.com/yyle88/yyle88"
	"github.com/yyle88/yyle88/internal/utils"
)

const username = "yyle88"

var organizationsSingleton []*yyle88.Organization
var onceFetchOrganizations sync.Once

func onceGetOrganizations() []*yyle88.Organization {
	onceFetchOrganizations.Do(func() {
		organizationsSingleton = done.VAE(yyle88.GetOrganizations(username)).Nice()
	})
	return organizationsSingleton
}

func TestGetOrganizations(t *testing.T) {
	t.Log(neatjsons.S(onceGetOrganizations()))
}

var mapOrganizationRepos = mutexmap.NewMap[string, []*yyle88.Repo](10)

func onceGetOrgRepos(organization *yyle88.Organization) []*yyle88.Repo {
	repos, _ := mapOrganizationRepos.Getset(organization.Name, func() []*yyle88.Repo {
		time.Sleep(time.Millisecond * 500)
		return rese.V1(yyle88.GetGithubRepos(organization.Name))
	})
	return repos
}

func TestFetchOrganizationRepos(t *testing.T) {
	organizations := onceGetOrganizations()
	require.NotEmpty(t, organizations)
	repos := onceGetOrgRepos(organizations[rand.IntN(len(organizations))])
	t.Log(neatjsons.S(repos))
}

type DocGenParam struct {
	shortName string
	titleLine string
}

func TestGenMarkdown(t *testing.T) {
	GenMarkdownTable(t, &DocGenParam{
		shortName: "README.md",
		titleLine: `| **<span style="font-size: 10px;">organization</span>** | **repo** |`,
	})
}

func TestGenMarkdownZhHans(t *testing.T) {
	GenMarkdownTable(t, &DocGenParam{
		shortName: "README.zh.md",
		titleLine: "| **组织** | **项目** |",
	})
}

func GenMarkdownTable(t *testing.T, arg *DocGenParam) {
	type orgRepo struct {
		orgName string
		repo    *yyle88.Repo
	}

	organizations := onceGetOrganizations()

	var results []*orgRepo
	var meaninglessRepos []*orgRepo
	for idx := 0; idx < 100; idx++ {
		var pieces = make([]*orgRepo, 0, len(organizations))
		for _, organization := range organizations {
			repos := onceGetOrgRepos(organization)

			if idx < len(repos) {
				if repo := repos[idx]; repo.Name == ".github" {
					meaninglessRepos = append(meaninglessRepos, &orgRepo{
						orgName: organization.Name,
						repo:    repo,
					})
				} else {
					pieces = append(pieces, &orgRepo{
						orgName: organization.Name,
						repo:    repo,
					})
				}
			}
		}
		rand.Shuffle(len(pieces), func(i, j int) {
			pieces[i], pieces[j] = pieces[j], pieces[i]
		})

		results = append(results, pieces...)
	}
	results = append(results, meaninglessRepos...)

	cardThemes := utils.GetRepoCardThemes()
	require.NotEmpty(t, cardThemes)

	rand.Shuffle(len(cardThemes), func(i, j int) {
		cardThemes[i], cardThemes[j] = cardThemes[j], cardThemes[i]
	})

	colors := utils.GetBadgeColors()
	require.NotEmpty(t, colors)

	rand.Shuffle(len(colors), func(i, j int) {
		colors[i], colors[j] = colors[j], colors[i]
	})

	ptx := utils.NewPTX()
	for _, organization := range organizations {
		ptx.Println(utils.MakeCustomSizeBadge(organization.Name, fmt.Sprintf("https://github.com/%s", organization.Name), colors[rand.IntN(len(colors))], 40, 125))
	}
	ptx.Println()

	ptx.Println(arg.titleLine)
	ptx.Println("|----------|----------|")

	for idx, one := range results {
		const templateLine = "[![Readme Card](https://github-readme-stats.vercel.app/api/pin/?username={{ username }}&repo={{ repo_name }}&theme={{ card_theme }}&unique={{ unique_uuid }})]({{ repo_link }})"

		rep := strings.NewReplacer(
			"{{ username }}", one.orgName,
			"{{ repo_name }}", one.repo.Name,
			"{{ card_theme }}", cardThemes[idx%len(cardThemes)],
			"{{ unique_uuid }}", uuid.New().String(),
			"{{ repo_link }}", one.repo.Link,
		)
		repoCardLink := rep.Replace(templateLine)

		orgBadgeLink := utils.MakeCustomSizeBadge(one.orgName, "https://github.com/"+one.orgName, colors[rand.IntN(len(colors))], 30, 80)

		ptx.Println(fmt.Sprintf("| %s | %s |", orgBadgeLink, repoCardLink))
	}

	stb := ptx.String()
	t.Log(stb)

	path := osmustexist.PATH(runpath.PARENT.Join(arg.shortName))
	t.Log(path)

	text := string(done.VAE(os.ReadFile(path)).Nice())
	t.Log(text)

	contentLines := strings.Split(text, "\n")
	sIdx := slices.Index(contentLines, "<!-- 这是一个注释，它不会在渲染时显示出来，这是组织项目列表的起始位置 -->")
	require.Positive(t, sIdx)
	eIdx := slices.Index(contentLines, "<!-- 这是一个注释，它不会在渲染时显示出来，这是组织项目列表的终止位置 -->")
	require.Positive(t, eIdx)

	require.Less(t, sIdx, eIdx)

	content := strings.Join(contentLines[:sIdx+1], "\n") + "\n" + "\n" +
		stb + "\n" +
		strings.Join(contentLines[eIdx:], "\n")
	t.Log(content)

	must.Done(os.WriteFile(path, []byte(content), 0666))
	t.Log("success")
}
