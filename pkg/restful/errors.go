package restful

import "github.com/golang/glog"

type Error struct {
	Code int32
	Msg  map[string]string // 支持多语言
}

type Errors struct {
	errors   map[int32]Error
	language string
}

func NewErrors() *Errors {
	return &Errors{
		errors:   make(map[int32]Error),
		language: "en",
	}
}

func (e *Errors) NewError(code int32, msg string) {
	err := Error{
		Code: code,
		Msg: map[string]string{
			e.language: msg,
		},
	}
	e.errors[code] = err
}

// 创建翻译
func (e *Errors) Translate(code int32, language string, msg string) {
	if _, ok := e.errors[code]; !ok {
		glog.Error("Error@Translate : code not exist")
		return
	}
	e.errors[code].Msg[language] = msg
}

// 获取错误消息
func (e *Errors) ErrorMsg(code int32) map[string]string {
	if _, ok := e.errors[code]; !ok {
		return nil
	}
	return e.errors[code].Msg
}
