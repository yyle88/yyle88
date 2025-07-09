package i18n_aboutmekeys_test

import (
	"testing"

	"github.com/yyle88/yyle88/locales/i18n_aboutmekeys"
)

func TestLoadI18nFiles(t *testing.T) {
	bundle, messageFiles := i18n_aboutmekeys.LoadI18nFiles()
	t.Log(len(messageFiles))
	t.Log(len(bundle.LanguageTags()))
}
