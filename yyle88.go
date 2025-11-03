package yyle88

import (
	"net/http"
	"os"
	"strings"
	"time"

	restyv2 "github.com/go-resty/resty/v2"
	"github.com/yyle88/erero"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/sortx"
	"github.com/yyle88/zaplog"
)

// Repo represents a GitHub repository with essential metadata
// Includes name, link, description, star count and last push time
//
// Repo 代表一个 GitHub 仓库及其基本元数据
// 包含名称、链接、描述、星标数和最后推送时间
type Repo struct {
	Name       string    `json:"name"`
	Link       string    `json:"html_url"`
	Desc       string    `json:"description"`
	Stargazers int       `json:"stargazers_count"`
	PushedAt   time.Time `json:"pushed_at"`
}

// GetGithubRepos fetches all public repositories for the given username
// Sorts repos by stars (descending) and recent activity
//
// GetGithubRepos 获取指定用户的所有公开仓库
// 按星标数（降序）和最近活跃度排序
func GetGithubRepos(username string) ([]*Repo, error) {
	var repos []*Repo

	response, err := newGithubRequest().SetPathParam("username", username).
		SetResult(&repos).
		Get("https://api.github.com/users/{username}/repos")
	if err != nil {
		return nil, erero.Wro(err)
	}
	if response.StatusCode() != http.StatusOK {
		return nil, erero.New(response.Status())
	}
	zaplog.SUG.Debugln(neatjsons.SxB(response.Body()))

	sortx.SortVStable(repos, func(a, b *Repo) bool {
		// 点开头的仓库排后面
		if strings.HasPrefix(a.Name, ".") || strings.HasPrefix(b.Name, ".") {
			return !strings.HasPrefix(a.Name, ".")
		}
		// 主页项目排后面，避免占据重要位置
		if a.Name == username || b.Name == username {
			return a.Name != username
		}
		// 星星数多的排前面
		if a.Stargazers != b.Stargazers {
			return a.Stargazers > b.Stargazers
		}
		// 星星数相同时，最近更新的排前面
		return a.PushedAt.After(b.PushedAt)
	})

	zaplog.SUG.Debugln(neatjsons.S(repos))
	return repos, nil
}

// newGithubRequest creates a new HTTP request with GitHub token if available
func newGithubRequest() *restyv2.Request {
	request := restyv2.New().SetTimeout(time.Minute).R()
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		request = request.SetHeader("Authorization", "token "+token)
	}
	return request
}

// Organization represents a GitHub organization
// Contains organization name and API endpoints
//
// Organization 代表一个 GitHub 组织
// 包含组织名称和 API 端点
type Organization struct {
	Name      string `json:"login"`     // Organization name // 组织名称
	Link      string `json:"url"`       // Organization URL // 组织链接
	ReposLink string `json:"repos_url"` // Repos API URL // 仓库接口链接
}

// GetOrganizations fetches all organizations that the user belongs to
// Returns list of organizations with their basic information
//
// GetOrganizations 获取用户所属的所有组织
// 返回包含基本信息的组织列表
func GetOrganizations(username string) ([]*Organization, error) {
	var organizations []*Organization

	response, err := newGithubRequest().SetPathParam("username", username).
		SetResult(&organizations).
		Get("https://api.github.com/users/{username}/orgs")
	if err != nil {
		return nil, erero.Wro(err)
	}
	if response.StatusCode() != http.StatusOK {
		return nil, erero.New(response.Status())
	}
	zaplog.SUG.Debugln(neatjsons.SxB(response.Body()))
	zaplog.SUG.Debugln(neatjsons.S(organizations))
	return organizations, nil
}

// GetOrganizationRepos fetches all repositories for the given organization
// Sorts repos with organization main repo first, then by stars
//
// GetOrganizationRepos 获取指定组织的所有仓库
// 组织主仓库排在最前，其余按星标数排序
func GetOrganizationRepos(orgName string) ([]*Repo, error) {
	var repos []*Repo

	response, err := newGithubRequest().SetPathParam("org", orgName).
		SetResult(&repos).
		Get("https://api.github.com/orgs/{org}/repos")
	if err != nil {
		return nil, erero.Wro(err)
	}
	if response.StatusCode() != http.StatusOK {
		return nil, erero.New(response.Status())
	}
	zaplog.SUG.Debugln(neatjsons.SxB(response.Body()))

	sortx.SortVStable(repos, func(a, b *Repo) bool {
		// 点开头的仓库排后面
		if strings.HasPrefix(a.Name, ".") || strings.HasPrefix(b.Name, ".") {
			return !strings.HasPrefix(a.Name, ".")
		}
		// 与组织同名的主项目排前面
		if a.Name == orgName || b.Name == orgName {
			return a.Name == orgName
		}
		// 星星数多的排前面
		if a.Stargazers != b.Stargazers {
			return a.Stargazers > b.Stargazers
		}
		// 星星数相同时，最近更新的排前面
		return a.PushedAt.After(b.PushedAt)
	})

	zaplog.SUG.Debugln(neatjsons.S(repos))
	return repos, nil
}
