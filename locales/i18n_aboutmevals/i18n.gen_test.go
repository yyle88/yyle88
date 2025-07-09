package i18n_aboutmevals_test

import (
	"testing"

	"github.com/yyle88/goi18n"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/runpath/runtestpath"
	"github.com/yyle88/syntaxgo"
	"github.com/yyle88/yyle88/locales/i18n_aboutmevals"
	"github.com/yyle88/zaplog"
)

func TestGenerate(t *testing.T) {
	bundle, messageFiles := i18n_aboutmevals.LoadI18nFiles()
	zaplog.SUG.Debugln(neatjsons.S(bundle.LanguageTags()))

	outputPath := osmustexist.FILE(runtestpath.SrcPath(t))
	options := goi18n.NewOptions().WithOutputPath(outputPath).WithPkgName(syntaxgo.GetPkgName(outputPath))
	t.Log(neatjsons.S(options))
	goi18n.Generate(messageFiles, options)
}
