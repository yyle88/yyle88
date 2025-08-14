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

func onceGetOrganizationRepos(organization *yyle88.Organization) []*yyle88.Repo {
	repos, _ := mapOrganizationRepos.Getset(organization.Name, func() []*yyle88.Repo {
		time.Sleep(time.Millisecond * 500)
		return rese.V1(yyle88.GetOrganizationRepos(organization.Name))
	})
	return repos
}

func TestFetchOrganizationRepos(t *testing.T) {
	organizations := onceGetOrganizations()
	require.NotEmpty(t, organizations)
	repos := onceGetOrganizationRepos(organizations[rand.IntN(len(organizations))])
	t.Log(neatjsons.S(repos))
	for _, repo := range repos {
		t.Log(repo.Name, repo.Stargazers)
	}
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
	const maxProjects = 10 // 限制最多显示10个项目

	for idx := 0; idx < 100 && len(results) < maxProjects; idx++ {
		var pieces = make([]*orgRepo, 0, len(organizations))
		for _, organization := range organizations {
			repos := onceGetOrganizationRepos(organization)

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

		// 只添加足够的项目达到maxProjects限制
		remainingSlots := maxProjects - len(results)
		if len(pieces) <= remainingSlots {
			results = append(results, pieces...)
		} else {
			results = append(results, pieces[:remainingSlots]...)
			break
		}
	}

	// 收集剩余的项目用于简洁展示
	var remainingRepos []*orgRepo
	for idx := 0; idx < 100; idx++ {
		for _, organization := range organizations {
			repos := onceGetOrganizationRepos(organization)
			if idx < len(repos) {
				if repo := repos[idx]; repo.Name != ".github" {
					// 检查是否已在results中
					found := false
					for _, result := range results {
						if result.orgName == organization.Name && result.repo.Name == repo.Name {
							found = true
							break
						}
					}
					if !found {
						remainingRepos = append(remainingRepos, &orgRepo{
							orgName: organization.Name,
							repo:    repo,
						})
					}
				}
			}
		}
	}

	// 添加meaninglessRepos到最后
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

	// 添加剩余项目的展示部分
	if len(remainingRepos) > 0 {
		ptx.Println()
		ptx.Println("---")
		ptx.Println()

		// 添加动态引言/格言
		if arg.shortName == "README.md" {
			ptx.Println("<div align=\"center\">")
			ptx.Println()
			ptx.Println("![Typing SVG](https://readme-typing-svg.herokuapp.com?font=Fira+Code&size=22&duration=4000&pause=1000&color=58A6FF&background=0D1117&center=true&vCenter=true&width=600&lines=🚀+Building+the+future%2C+one+commit+at+a+time;💡+Innovation+through+elegant+code;🌟+Turning+ideas+into+reality)")
			ptx.Println()
			ptx.Println("*\"Code is like humor. When you have to explain it, it's bad.\"* – Cory House")
			ptx.Println()
			ptx.Println("</div>")
		} else {
			ptx.Println("<div align=\"center\">")
			ptx.Println()
			ptx.Println("![Typing SVG](https://readme-typing-svg.herokuapp.com?font=Fira+Code&size=22&duration=4000&pause=1000&color=58A6FF&background=0D1117&center=true&vCenter=true&width=600&lines=🚀+一行代码改变世界;💡+用优雅的代码创新未来;🌟+将想法变为现实)")
			ptx.Println()
			ptx.Println("*\"优雅的代码是可以自我解释的代码\"* – 代码之道")
			ptx.Println()
			ptx.Println("</div>")
		}

		ptx.Println()
		ptx.Println("---")
		ptx.Println()

		// 统计信息
		totalStars := 0
		totalRepos := len(results) + len(remainingRepos)
		orgStarMap := make(map[string]int)

		// 统计前10个项目的stars
		for _, one := range results {
			if one.repo.Name != ".github" {
				totalStars += one.repo.Stargazers
				orgStarMap[one.orgName] += one.repo.Stargazers
			}
		}

		// 统计剩余项目的stars
		for _, one := range remainingRepos {
			totalStars += one.repo.Stargazers
			orgStarMap[one.orgName] += one.repo.Stargazers
		}

		// 添加增强版统计徽章
		if arg.shortName == "README.md" {
			ptx.Println("<div align=\"center\">")
			ptx.Println()

			// 基础统计 - 使用更炫酷的颜色
			ptx.Printf("![Total Stars](https://img.shields.io/badge/⭐_Total_Stars-%d-FFD700?style=for-the-badge&logo=github&logoColor=white&labelColor=FF6B6B)\n", totalStars)
			ptx.Printf("![Total Repos](https://img.shields.io/badge/📁_Total_Repos-%d-4ECDC4?style=for-the-badge&logo=git&logoColor=white&labelColor=45B7D1)\n", totalRepos)
			ptx.Printf("![Organizations](https://img.shields.io/badge/🏢_Organizations-%d-96CEB4?style=for-the-badge&logo=organization&logoColor=white&labelColor=FFEAA7)\n", len(organizations))

			ptx.Println()

			// 添加访客计数和年限统计
			ptx.Printf("![Profile Views](https://komarev.com/ghpvc/?username=yyle88&style=for-the-badge&color=blueviolet&label=PROFILE+VIEWS)\n")
			ptx.Printf("![Years Badge](https://badges.pufler.dev/years/yyle88?style=for-the-badge&color=blue&logo=github)\n")
			ptx.Printf("![Repos Badge](https://badges.pufler.dev/repos/yyle88?style=for-the-badge&color=success&logo=github)\n")

			ptx.Println()
			ptx.Println("</div>")
		} else {
			ptx.Println("<div align=\"center\">")
			ptx.Println()

			ptx.Printf("![总Stars数](https://img.shields.io/badge/⭐_总Stars数-%d-FFD700?style=for-the-badge&logo=github&logoColor=white&labelColor=FF6B6B)\n", totalStars)
			ptx.Printf("![总项目数](https://img.shields.io/badge/📁_总项目数-%d-4ECDC4?style=for-the-badge&logo=git&logoColor=white&labelColor=45B7D1)\n", totalRepos)
			ptx.Printf("![组织数量](https://img.shields.io/badge/🏢_组织数量-%d-96CEB4?style=for-the-badge&logo=organization&logoColor=white&labelColor=FFEAA7)\n", len(organizations))

			ptx.Println()

			ptx.Printf("![访问量](https://komarev.com/ghpvc/?username=yyle88&style=for-the-badge&color=blueviolet&label=访问量)\n")
			ptx.Printf("![编程年限](https://badges.pufler.dev/years/yyle88?style=for-the-badge&color=blue&logo=github)\n")
			ptx.Printf("![仓库总数](https://badges.pufler.dev/repos/yyle88?style=for-the-badge&color=success&logo=github)\n")

			ptx.Println()
			ptx.Println("</div>")
		}

		ptx.Println()
		ptx.Println("---")
		ptx.Println()

		// 显示其他项目的简洁列表
		if arg.shortName == "README.md" {
			ptx.Println("<h3 align=\"center\">🚀 More Projects</h3>")
		} else {
			ptx.Println("<h3 align=\"center\">🚀 更多项目</h3>")
		}
		ptx.Println()
		ptx.Println("<div align=\"center\">")
		ptx.Println()

		// 按组织分组显示剩余项目
		orgProjects := make(map[string][]*orgRepo)
		for _, repo := range remainingRepos {
			orgProjects[repo.orgName] = append(orgProjects[repo.orgName], repo)
		}

		for _, organization := range organizations {
			if projects, exists := orgProjects[organization.Name]; exists && len(projects) > 0 {
				ptx.Printf("**%s** • ", strings.ToUpper(organization.Name))
				for i, project := range projects {
					if i > 0 {
						ptx.Print(" • ")
					}
					ptx.Printf("[%s](https://github.com/%s/%s)", project.repo.Name, project.orgName, project.repo.Name)
					if project.repo.Stargazers > 0 {
						ptx.Printf("⭐%d", project.repo.Stargazers)
					}
				}
				ptx.Println()
				ptx.Println()
			}
		}

		ptx.Println("</div>")
		ptx.Println()
		ptx.Println("---")
		ptx.Println()

		// 添加超级酷炫的技术栈展示
		if arg.shortName == "README.md" {
			ptx.Println("<h3 align=\"center\">🛠️ Tech Arsenal & Skills</h3>")
		} else {
			ptx.Println("<h3 align=\"center\">🛠️ 技术武器库</h3>")
		}
		ptx.Println()
		ptx.Println("<div align=\"center\">")
		ptx.Println()

		// 分类展示技术栈
		if arg.shortName == "README.md" {
			ptx.Println("### 🚀 **Languages & Frameworks**")
		} else {
			ptx.Println("### 🚀 **编程语言与框架**")
		}
		ptx.Println()

		// 主要编程语言
		mainTechStacks := []string{
			"![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white&labelColor=E10098)",
			"![Kratos](https://img.shields.io/badge/Kratos-7C3AED?style=for-the-badge&logo=go-kratos&logoColor=white&labelColor=FF6B6B)",
			"![Gin](https://img.shields.io/badge/Gin-00ADD8?style=for-the-badge&logo=gin&logoColor=white&labelColor=4ECDC4)",
			"![GORM](https://img.shields.io/badge/GORM-00D9FF?style=for-the-badge&logo=go&logoColor=white&labelColor=95DE64)",
		}

		for _, tech := range mainTechStacks {
			ptx.Print(tech + " ")
		}
		ptx.Println()
		ptx.Println()

		if arg.shortName == "README.md" {
			ptx.Println("### 🔧 **DevOps & Infrastructure**")
		} else {
			ptx.Println("### 🔧 **运维与基础设施**")
		}
		ptx.Println()

		devopsTechStacks := []string{
			"![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white&labelColor=FF6B35)",
			"![Kubernetes](https://img.shields.io/badge/Kubernetes-326CE5?style=for-the-badge&logo=kubernetes&logoColor=white&labelColor=7209B7)",
			"![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?style=for-the-badge&logo=github-actions&logoColor=white&labelColor=FF6347)",
		}

		for _, tech := range devopsTechStacks {
			ptx.Print(tech + " ")
		}
		ptx.Println()
		ptx.Println()

		if arg.shortName == "README.md" {
			ptx.Println("### 💾 **Databases & Message Queues**")
		} else {
			ptx.Println("### 💾 **数据库与消息队列**")
		}
		ptx.Println()

		dbTechStacks := []string{
			"![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white&labelColor=FF4081)",
			"![MongoDB](https://img.shields.io/badge/MongoDB-4EA94B?style=for-the-badge&logo=mongodb&logoColor=white&labelColor=FFA726)",
			"![Redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white&labelColor=9C27B0)",
			"![Kafka](https://img.shields.io/badge/Apache_Kafka-231F20?style=for-the-badge&logo=apache-kafka&logoColor=white&labelColor=00BCD4)",
		}

		for _, tech := range dbTechStacks {
			ptx.Print(tech + " ")
		}
		ptx.Println()
		ptx.Println()

		// 添加技能进度展示
		if arg.shortName == "README.md" {
			ptx.Println("### ⚡ **Skill Levels**")
			ptx.Println()
			ptx.Println("```text")
			ptx.Println("Go Programming    ████████████████████   100%")
			ptx.Println("Microservices     ██████████████████░░    90%")
			ptx.Println("Docker/K8s        ████████████████░░░░    80%")
			ptx.Println("System Design     █████████████░░░░░░░    65%")
			ptx.Println("Cloud Architecture ██████████████░░░░░░    70%")
			ptx.Println("```")
		} else {
			ptx.Println("### ⚡ **技能等级**")
			ptx.Println()
			ptx.Println("```text")
			ptx.Println("Go 编程          ████████████████████   100%")
			ptx.Println("微服务架构        ██████████████████░░    90%")
			ptx.Println("容器化部署        ████████████████░░░░    80%")
			ptx.Println("系统设计          █████████████░░░░░░░    65%")
			ptx.Println("云架构设计        ██████████████░░░░░░    70%")
			ptx.Println("```")
		}

		ptx.Println()
		ptx.Println("</div>")
		ptx.Println()
		ptx.Println("---")
		ptx.Println()

		// 添加GitHub统计
		if arg.shortName == "README.md" {
			ptx.Println("<h3 align=\"center\">📊 GitHub Stats</h3>")
		} else {
			ptx.Println("<h3 align=\"center\">📊 GitHub 统计</h3>")
		}
		ptx.Println()
		ptx.Println("<div align=\"center\">")
		ptx.Println()
		ptx.Println("![GitHub Stats](https://github-readme-stats.vercel.app/api?username=yyle88&show_icons=true&theme=radical)")
		ptx.Println()
		ptx.Println("![Top Languages](https://github-readme-stats.vercel.app/api/top-langs/?username=yyle88&layout=compact&theme=radical)")
		ptx.Println()
		ptx.Println("</div>")
		ptx.Println()

		// 添加活动图表
		ptx.Println("---")
		ptx.Println()
		if arg.shortName == "README.md" {
			ptx.Println("<h3 align=\"center\">📈 Activity Graph</h3>")
		} else {
			ptx.Println("<h3 align=\"center\">📈 活动图表</h3>")
		}
		ptx.Println()
		ptx.Println("<div align=\"center\">")
		ptx.Println()
		ptx.Println("![Activity Graph](https://github-readme-activity-graph.vercel.app/graph?username=yyle88&theme=react-dark)")
		ptx.Println()
		ptx.Println("</div>")
		ptx.Println()
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
