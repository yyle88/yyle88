package yyle88

import (
	"net/http"
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
	response, err := restyv2.New().SetTimeout(time.Minute).R().
		SetPathParam("username", username).
		SetResult(&repos).
		Get("https://api.github.com/users/{username}/repos")
	if err != nil {
		return nil, erero.Wro(err)
	}
	if response.StatusCode() != http.StatusOK {
		return nil, erero.New(response.Status())
	}

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
	return repos, nil
}
