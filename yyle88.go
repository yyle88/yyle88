package yyle88

import (
	"net/http"
	"os"
	"strings"
	"time"

	restyv2 "github.com/go-resty/resty/v2"
	"github.com/yyle88/erero"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/sortslice"
	"github.com/yyle88/zaplog"
)

type Repo struct {
	Name       string    `json:"name"`
	Link       string    `json:"html_url"`
	Desc       string    `json:"description"`
	Stargazers int       `json:"stargazers_count"`
	PushedAt   time.Time `json:"pushed_at"`
}

func GetGithubRepos(username string) ([]*Repo, error) {
	var repos []*Repo

	// 从环境变量读取 GitHub Token
	githubToken := os.Getenv("GITHUB_TOKEN")

	// 使用 Token 添加 Authorization 请求头
	request := restyv2.New().SetTimeout(time.Minute).R()
	if githubToken != "" {
		request = request.SetHeader("Authorization", "token "+githubToken)
	}
	response, err := request.SetPathParam("username", username).
		SetResult(&repos).
		Get("https://api.github.com/users/{username}/repos")
	if err != nil {
		return nil, erero.Wro(err)
	}
	if response.StatusCode() != http.StatusOK {
		return nil, erero.New(response.Status())
	}
	zaplog.SUG.Debugln(neatjsons.SxB(response.Body()))

	sortslice.SortVStable(repos, func(a, b *Repo) bool {
		if strings.HasPrefix(a.Name, ".") || strings.HasPrefix(b.Name, ".") {
			return !strings.HasPrefix(a.Name, ".")
		} else if a.Name == username || b.Name == username {
			return a.Name != username //当是主页项目时把它排在最后面，避免排的太靠前占据重要的位置
		} else if a.Stargazers > b.Stargazers {
			return true //星多者排前面
		} else if a.Stargazers < b.Stargazers {
			return false //星少者排后面
		} else {
			return a.PushedAt.After(b.PushedAt) //同样星星时最近有更新的排前面
		}
	})

	zaplog.SUG.Debugln(neatjsons.S(repos))
	return repos, nil
}

type Organization struct {
	Name      string `json:"login"`     // 组织名称
	Link      string `json:"url"`       // 组织链接
	ReposLink string `json:"repos_url"` // 组织链接
}

func GetOrganizations(username string) ([]*Organization, error) {
	var organizations []*Organization

	// 从环境变量读取 GitHub Token
	githubToken := os.Getenv("GITHUB_TOKEN")

	// 使用 Token 添加 Authorization 请求头
	request := restyv2.New().SetTimeout(time.Minute).R()
	if githubToken != "" {
		request = request.SetHeader("Authorization", "Bearer "+githubToken)
	}

	// 请求获取用户的组织信息
	response, err := request.SetPathParam("username", username).
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
