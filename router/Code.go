package router

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"reflect"
	_ "github.com/davyxu/cellnet/codec/httpform"
	_ "github.com/davyxu/cellnet/codec/json"
)

type HttpEmailCode struct {
	UserEmail string
	vCode string
}
func init() {
	cellnet.RegisterHttpMeta(&cellnet.HttpMeta{
		URL:          "/register",
		Method:       "POST",
		RequestCodec: codec.MustGetCodec("httpform"),
		RequestType:  reflect.TypeOf((*HttpRegisterREQ)(nil)).Elem(),
		ResponseCodec: codec.MustGetCodec("json"),
		ResponseType:  reflect.TypeOf((*HttpTokenACK)(nil)).Elem(),
	})
}