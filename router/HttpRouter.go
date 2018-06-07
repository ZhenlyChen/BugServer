package router

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"reflect"
	"fmt"
	_ "github.com/davyxu/cellnet/codec/httpform"
	_ "github.com/davyxu/cellnet/codec/json"
)


type HttpEchoREQ struct {
	UserName string
}

type HttpEchoACK struct {
	Token  string
	Status int32
}

func (self *HttpEchoREQ) String() string { return fmt.Sprintf("%+v", *self) }
func (self *HttpEchoACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {

	cellnet.RegisterHttpMeta(&cellnet.HttpMeta{
		URL:          "/hello",
		Method:       "GET",
		RequestCodec: codec.MustGetCodec("httpform"),
		RequestType:  reflect.TypeOf((*HttpEchoREQ)(nil)).Elem(),

		ResponseCodec: codec.MustGetCodec("json"),
		ResponseType:  reflect.TypeOf((*HttpEchoACK)(nil)).Elem(),
	})
}