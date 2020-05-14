package restful

import (
	"sync"

	"github.com/golang/glog"
)

type Error struct {
	Code string
	Msg  map[string]string // 支持多语言
}

type Errors struct {
	errors   map[string]Error
	language string
}

type ErrorsBucket struct {
	errorsMap []*Errors
}

var errorBucketInsObj *ErrorsBucket
var errorBucketOnce sync.Once

func errorBucketIns() *ErrorsBucket {
	errorBucketOnce.Do(func() {
		errorBucketInsObj = &ErrorsBucket{
			errorsMap: make([]*Errors, 0),
		}
	})
	return errorBucketInsObj
}

func NewErrors() *Errors {
	errors := &Errors{
		errors:   make(map[string]Error),
		language: "en",
	}
	errorBucketIns().errorsMap = append(errorBucketIns().errorsMap, errors)
	return errors
}

func (e *Errors) NewError(code string, msg string) {
	err := Error{
		Code: code,
		Msg: map[string]string{
			e.language: msg,
		},
	}
	e.errors[code] = err
}

// 创建翻译
func (e *Errors) Translate(code string, language string, msg string) {
	if _, ok := e.errors[code]; !ok {
		glog.Error("Error@Translate : code not exist")
		return
	}
	e.errors[code].Msg[language] = msg
}

// 获取错误消息
func (e *Errors) ErrorMsg(code string) map[string]string {
	if _, ok := e.errors[code]; !ok {
		return nil
	}
	return e.errors[code].Msg
}
