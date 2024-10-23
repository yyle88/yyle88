package yyle88_test

import (
	"os"
	"testing"

	"github.com/yyle88/done"
	"github.com/yyle88/must"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/runpath/runtestpath"
)

func TestExample(t *testing.T) {
	path := runtestpath.SrcPath(t)

	t.Log(path)

	must.True(osmustexist.IsFile(path))

	data := done.VAE(os.ReadFile(path)).Nice()

	t.Log(string(data))
}
