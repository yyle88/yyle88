package i18n_message_test

import (
	"testing"

	"github.com/yyle88/yyle88/locales/i18n_message"
)

func TestLoadI18nFiles(t *testing.T) {
	bundle, messageFiles := i18n_message.LoadI18nFiles()
	t.Log(len(messageFiles))
	t.Log(len(bundle.LanguageTags()))
}
