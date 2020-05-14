package restful

import (
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
)

type HandlerOpts struct {
	ParseFunc     func(map[string][]string, []byte) (map[string]interface{}, error)
	MakeErrorFunc func(*Errors)
	Filters       []Filter
}

func NewHandler(handler RestHandler, opt HandlerOpts) func(http.ResponseWriter, *http.Request) {
	errosMgr := NewErrors()
	// 制作错误数据
	if opt.MakeErrorFunc != nil {
		opt.MakeErrorFunc(errosMgr)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err    error
			body   []byte
			params map[string]interface{}
			resp   *Response
		)
		// 创建应答数据
		resp = &Response{
			errMgr:  errosMgr,
			isWrite: true, // 默认自动向客户端写入应答数据
		}
		defer r.Body.Close()
		body, err = ioutil.ReadAll(r.Body)
		//if len(body) > 0 {
		//	glog.Info(string(body))
		//}
		if err != nil {
			glog.Error(err)
			resp.Error("", err.Error()).Write(w)
			return
		}

		// 调用解析器
		if opt.ParseFunc != nil {
			params, err = opt.ParseFunc(r.URL.Query(), body)
			if err != nil {
				glog.Error(err)
				resp.Error("", err.Error()).Write(w)
				return
			}
		}
		// 调用过滤器
		for _, f := range opt.Filters {
			params, err = f.Processor(r, params)
			if err != nil {
				glog.Error(err)
				resp.Error("", err.Error()).Write(w)
				return // 过滤失败,不调用handler处理
			}
		}
		handler(w, r, params, resp)
		if resp.isWrite {
			resp.Write(w)
		}
	}
}
