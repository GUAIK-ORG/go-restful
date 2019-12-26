package restful

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Response struct {
	Status    int32                  `json:"status"`     // 0: 成功 -1: 失败
	ErrorCode int32                  `json:"error_code"` // 错误码
	ErrorMsg  map[string]string      `json:"error_msg"`  // 错误信息
	Body      map[string]interface{} `json:"body"`       // 返回数据
	errMgr    *Errors                // 错误管理器
	isWrite   bool                   // 是否向客户端写入应答数据
}

func (r *Response) Success(body map[string]interface{}) *Response {
	r.Status = 0
	r.ErrorCode = 0
	r.ErrorMsg = nil
	r.Body = body
	return r
}

// 设置是否自动向客户端写入应答数据
// 如果需要自行处理的返回数据的话可以设置为false
func (r *Response) UseWrite(b bool) {
	r.isWrite = b
}

func (r *Response) Error(errCode int32, errMsg string) *Response {
	r.Status = -1
	r.ErrorCode = errCode
	r.ErrorMsg = map[string]string{"en": errMsg}
	return r
}

func (r *Response) UseError(errCode int32) *Response {
	r.Status = -1
	r.ErrorCode = errCode
	r.ErrorMsg = r.errMgr.ErrorMsg(errCode)
	return r
}

// 向客户端写入响应个数据
func (r *Response) Write(w http.ResponseWriter) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(r)
	// glog.Info(string(bf.String()))
	w.Write(bf.Bytes())
}
