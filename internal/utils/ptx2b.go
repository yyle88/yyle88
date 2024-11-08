package utils

import (
	"bytes"
	"fmt"

	"github.com/yyle88/done"
)

//goland:noinspection GoExportedFuncWithUnexportedType
func NewPTX() *print2Bytes {
	return &print2Bytes{}
}

type print2Bytes struct{ bytes.Buffer }

func (T *print2Bytes) Println(ax ...interface{}) (n int) {
	return done.VE(fmt.Fprintln(T, ax...)).Done()
}
