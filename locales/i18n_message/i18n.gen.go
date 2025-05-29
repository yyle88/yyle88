package i18n_message

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func NewRepoTableRepoTitle() string {
	return "REPO_TABLE_REPO_TITLE"
}

func I18nRepoTableRepoTitle() *i18n.LocalizeConfig {
	return &i18n.LocalizeConfig{
		MessageID: "REPO_TABLE_REPO_TITLE",
	}
}

func NewRepoTableTitleDesc() string {
	return "REPO_TABLE_TITLE_DESC"
}

func I18nRepoTableTitleDesc() *i18n.LocalizeConfig {
	return &i18n.LocalizeConfig{
		MessageID: "REPO_TABLE_TITLE_DESC",
	}
}

func NewRepoTableTitleName() string {
	return "REPO_TABLE_TITLE_NAME"
}

func I18nRepoTableTitleName() *i18n.LocalizeConfig {
	return &i18n.LocalizeConfig{
		MessageID: "REPO_TABLE_TITLE_NAME",
	}
}
