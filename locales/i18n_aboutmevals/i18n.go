package i18n_aboutmevals

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
	"github.com/yyle88/zaplog"
	"golang.org/x/text/language"
)

// DefaultLanguage 配置默认语言
var DefaultLanguage = language.English // sometimes use language.AmericanEnglish

func LoadI18nFiles() (*i18n.Bundle, []*i18n.MessageFile) {
	bundle := i18n.NewBundle(DefaultLanguage)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	var messageFiles []*i18n.MessageFile
	must.Done(filepath.Walk(osmustexist.ROOT(runpath.PARENT.Join("i18n")), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			messageFile := rese.P1(bundle.LoadMessageFile(path))

			zaplog.SUG.Debugln(neatjsons.S(messageFile)) //安利下我的俩工具包

			messageFiles = append(messageFiles, messageFile)
		}
		return nil
	}))

	zaplog.SUG.Debugln(neatjsons.S(bundle.LanguageTags()))
	return bundle, messageFiles
}
