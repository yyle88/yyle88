package i18n_message

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func I18nRepoTableRepoTitle() *i18n.LocalizeConfig {
	const messageID = "REPO_TABLE_REPO_TITLE"
	return &i18n.LocalizeConfig{
		MessageID: messageID,
	}
}

func I18nRepoTableTitleDesc() *i18n.LocalizeConfig {
	const messageID = "REPO_TABLE_TITLE_DESC"
	return &i18n.LocalizeConfig{
		MessageID: messageID,
	}
}

func I18nRepoTableTitleName() *i18n.LocalizeConfig {
	const messageID = "REPO_TABLE_TITLE_NAME"
	return &i18n.LocalizeConfig{
		MessageID: messageID,
	}
}
