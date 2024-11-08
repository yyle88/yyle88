package yyle88

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/yyle88/done"
	"github.com/yyle88/must"
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

func MustGetGithubRepos(username string) []*Repo {
	var repos []*Repo
	resp := done.VPE(resty.New().SetTimeout(time.Minute).R().
		SetPathParam("username", username).
		SetResult(&repos).
		Get("https://api.github.com/users/{username}/repos")).Nice()
	must.Equals(resp.StatusCode(), http.StatusOK)

	sortslice.SortVStable(repos, func(a, b *Repo) bool {
		if a.Stargazers > b.Stargazers {
			return true
		} else if a.Stargazers < b.Stargazers {
			return false
		} else {
			return a.PushedAt.After(b.PushedAt)
		}
	})
	zaplog.SUG.Info(neatjsons.S(repos))
	return repos
}
