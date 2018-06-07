package router

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	_ "github.com/davyxu/cellnet/codec/httpform"
	_ "github.com/davyxu/cellnet/codec/json"
	"reflect"
)

type HttpLoginREQ struct {
	UserName string
	Password string
}

func init() {
	cellnet.RegisterHttpMeta(&cellnet.HttpMeta{
		URL:          "/login",
		Method:       "POST",
		RequestCodec: codec.MustGetCodec("httpform"),
		RequestType:  reflect.TypeOf((*HttpLoginREQ)(nil)).Elem(),
		ResponseCodec: codec.MustGetCodec("json"),
		ResponseType:  reflect.TypeOf((*HttpTokenACK)(nil)).Elem(),
	})
}