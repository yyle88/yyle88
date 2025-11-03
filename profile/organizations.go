package profile

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/yyle88/mutexmap"
	"github.com/yyle88/rese"
	"github.com/yyle88/yyle88"
	"github.com/yyle88/yyle88/internal/utils"
)

const username = "yyle88"

var (
	organizationsSingleton []*yyle88.Organization
	onceFetchOrganizations sync.Once
)

// getOrganizationsCached fetches organizations with singleton pattern
// Ensures organizations are loaded only once per session
//
// getOrganizationsCached 使用单例模式获取组织列表
// 确保每个会话只加载一次组织数据
func getOrganizationsCached() []*yyle88.Organization {
	onceFetchOrganizations.Do(func() {
		organizationsSingleton = rese.V1(yyle88.GetOrganizations(username))
	})
	return organizationsSingleton
}

// organizationReposCache stores fetched organization repositories
// Prevents redundant API calls for the same organization
//
// organizationReposCache 存储已获取的组织仓库
// 防止对同一组织进行冗余的 API 调用
var organizationReposCache = mutexmap.NewMap[string, []*yyle88.Repo](10)

// getOrganizationReposCached fetches organization repos with caching
// Adds delay to respect API rate limits
//
// getOrganizationReposCached 使用缓存获取组织仓库
// 添加延迟以遵守 API 速率限制
func getOrganizationReposCached(organization *yyle88.Organization) []*yyle88.Repo {
	repos, _ := organizationReposCache.Getset(organization.Name, func() []*yyle88.Repo {
		time.Sleep(time.Millisecond * 500)
		return rese.V1(yyle88.GetOrganizationRepos(organization.Name))
	})
	return repos
}

// organizationRepo represents a repository belonging to an organization
// Combines organization name with repository details
//
// organizationRepo 代表属于某个组织的仓库
// 组合了组织名称和仓库详情
type organizationRepo struct {
	orgName string
	repo    *yyle88.Repo
}

// isGithubRepo checks if this is a .github configuration repository
// isGithubRepo 检查是否为 .github 配置仓库
func (or *organizationRepo) isGithubRepo() bool {
	return or.repo.Name == ".github"
}

// repoCollection organizes repositories from multiple organizations
// Separates repos into featured, remaining, and meta categories
//
// repoCollection 组织来自多个组织的仓库
// 将仓库分为精选、剩余和元数据类别
type repoCollection struct {
	organizations  []*yyle88.Organization // All organizations // 所有组织
	featuredRepos  []*organizationRepo    // Top repos to display as cards // 以卡片展示的顶级仓库
	remainingRepos []*organizationRepo    // Other repos in table // 表格中的其他仓库
	metaRepos      []*organizationRepo    // Config repos like .github // 配置仓库如 .github
}

// newRepoCollection creates a new repository collection
// Initializes with cached organizations
//
// newRepoCollection 创建新的仓库集合
// 使用缓存的组织初始化
func newRepoCollection() *repoCollection {
	return &repoCollection{
		organizations: getOrganizationsCached(),
	}
}

// collectRepos gathers repos using round-robin pattern from organizations
// Each round takes one repo from each org, ensuring fair distribution
// Separates .github repos and sorts remaining ones by stars
//
// collectRepos 使用轮询模式从组织收集仓库
// 每轮从每个组织取一个仓库，确保公平分配
// 分离 .github 仓库并按星标数排序剩余仓库
func (rc *repoCollection) collectRepos() {
	const maxFeaturedProjects = 10
	const maxRounds = 100 // 最多收集100轮，足够覆盖所有项目
	var allRepos []*organizationRepo

	// 按轮次收集：每个组织轮流拿第0个、第1个...（实现"每个组织最高星→次高星"的效果）
	for round := 0; round < maxRounds; round++ {
		var roundRepos []*organizationRepo
		for _, organization := range rc.organizations {
			repos := getOrganizationReposCached(organization)
			if round < len(repos) {
				repo := repos[round]
				orgRepo := &organizationRepo{
					orgName: organization.Name,
					repo:    repo,
				}
				// 在收集时就区分 .github 项目
				if orgRepo.isGithubRepo() {
					rc.metaRepos = append(rc.metaRepos, orgRepo)
				} else {
					roundRepos = append(roundRepos, orgRepo)
				}
			}
		}
		// 同一轮内随机打乱
		rand.Shuffle(len(roundRepos), func(i, j int) {
			roundRepos[i], roundRepos[j] = roundRepos[j], roundRepos[i]
		})
		allRepos = append(allRepos, roundRepos...)
	}

	// 选择 featured repos 和 remaining repos
	if len(allRepos) <= maxFeaturedProjects {
		rc.featuredRepos = allRepos
	} else {
		rc.featuredRepos = allRepos[:maxFeaturedProjects]
		rc.remainingRepos = allRepos[maxFeaturedProjects:]

		// 剩余项目按星星数降序排序
		slices.SortFunc(rc.remainingRepos, func(a, b *organizationRepo) int {
			return b.repo.Stargazers - a.repo.Stargazers
		})
	}
}

// profileMarkdown represents markdown content generation for GitHub profile
// Contains repo collection and styling configuration
//
// profileMarkdown 代表 GitHub 个人资料的 markdown 内容生成
// 包含仓库集合和样式配置
type profileMarkdown struct {
	collection *repoCollection // Repository collection // 仓库集合
	cardThemes []string        // Available card themes // 可用的卡片主题
	language   string          // Language code // 语言代码
}

// newProfileMarkdown creates a profile markdown generator with styling
// Collects repos and shuffles themes/colors for variety
//
// newProfileMarkdown 创建带样式的个人资料 markdown 生成器
// 收集仓库并打乱主题/颜色以增加多样性
func newProfileMarkdown(language string) *profileMarkdown {
	collection := newRepoCollection()
	collection.collectRepos()

	cardThemes := utils.GetRepoCardThemes()
	rand.Shuffle(len(cardThemes), func(i, j int) {
		cardThemes[i], cardThemes[j] = cardThemes[j], cardThemes[i]
	})

	return &profileMarkdown{
		collection: collection,
		cardThemes: cardThemes,
		language:   language,
	}
}

// generateOrgsBadges creates badge links for all organizations
// Each badge is a clickable link to the organization page
//
// generateOrgsBadges 为所有组织创建徽章链接
// 每个徽章都是指向组织页面的可点击链接
func (mc *profileMarkdown) generateOrgsBadges() string {
	ptx := utils.NewPTX()
	colors := newShuffleBadgeColors()
	for _, organization := range mc.collection.organizations {
		ptx.Println(utils.MakeCustomSizeBadge(
			organization.Name,
			fmt.Sprintf("https://github.com/%s", organization.Name),
			colors[rand.IntN(len(colors))],
			40, 125,
		))
	}
	ptx.Println()
	return ptx.String()
}

func newShuffleBadgeColors() []string {
	colors := utils.GetBadgeColors()
	rand.Shuffle(len(colors), func(i, j int) {
		colors[i], colors[j] = colors[j], colors[i]
	})
	return colors
}

// generateProjectsTable creates table of featured organization projects
// Shows top projects with repository cards and organization badges
//
// generateProjectsTable 创建精选组织项目表格
// 展示顶级项目及其仓库卡片和组织徽章
func (mc *profileMarkdown) generateProjectsTable(titleLine string) string {
	ptx := utils.NewPTX()
	ptx.Println(titleLine)
	ptx.Println("|----------|----------|")

	colors := newShuffleBadgeColors()
	for idx, item := range mc.collection.featuredRepos {
		const templateLine = "[![Readme Card](https://github-readme-stats.vercel.app/api/pin/?username={{ username }}&repo={{ repo_name }}&theme={{ card_theme }}&unique={{ unique_uuid }})]({{ repo_link }})"
		rep := strings.NewReplacer(
			"{{ username }}", item.orgName,
			"{{ repo_name }}", item.repo.Name,
			"{{ card_theme }}", mc.cardThemes[idx%len(mc.cardThemes)],
			"{{ unique_uuid }}", uuid.New().String(),
			"{{ repo_link }}", item.repo.Link,
		)
		repoCardLink := rep.Replace(templateLine)
		orgBadgeLink := utils.MakeCustomSizeBadge(
			item.orgName,
			"https://github.com/"+item.orgName,
			colors[rand.IntN(len(colors))],
			30, 80,
		)
		ptx.Println(fmt.Sprintf("| %s | %s |", orgBadgeLink, repoCardLink))
	}

	ptx.Println()
	return ptx.String()
}

// generateMoreProjects creates 4-column table of remaining projects
// Includes both regular repos and .github config repos
// Config repos are placed at the end of the table
//
// generateMoreProjects 创建剩余项目的 4 列表格
// 包含常规仓库和 .github 配置仓库
// 配置仓库放在表格末尾
func (mc *profileMarkdown) generateMoreProjects() string {
	ptx := utils.NewPTX()
	ptx.Println("<div>")
	ptx.Println()

	if mc.language == "en" {
		ptx.Println("| Repo | Repo | Repo | Repo |")
	} else {
		ptx.Println("| 项目 | 项目 | 项目 | 项目 |")
	}
	ptx.Println("| :--: | :--: | :--: | :--: |")

	// 将所有剩余项目收集到一个列表（把 .github 项目放到末尾）
	allRepos := append([]*organizationRepo{}, mc.collection.remainingRepos...)
	allRepos = append(allRepos, mc.collection.metaRepos...)

	colors := newShuffleBadgeColors()

	// 按4列输出
	const columns = 4
	for i := 0; i < len(allRepos); i += columns {
		ptx.Print("|")
		for j := 0; j < columns; j++ {
			idx := i + j
			if idx < len(allRepos) {
				project := allRepos[idx]
				badge := utils.MakeCustomSizeBadge(
					project.repo.Name,
					fmt.Sprintf("https://github.com/%s/%s", project.orgName, project.repo.Name),
					colors[rand.IntN(len(colors))],
					24, 0, // height=24, width=auto
				)
				ptx.Print(badge)
			}
			ptx.Print(" | ")
		}
		ptx.Println()
	}

	ptx.Println()
	ptx.Println("</div>")
	return ptx.String()
}
