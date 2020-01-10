package restful

import (
	"net/http"
)

type SchedulerOpt struct {
	UseCORS     bool
	AllowOrigin string
}

type Scheduler struct {
	// 处理函数
	PostHandler   HttpHandler // 创建
	GetHandler    HttpHandler // 获取
	PutHandler    HttpHandler // 更新（全部更新）
	PatchHandler  HttpHandler // 更新（部分更新）
	DeleteHandler HttpHandler // 删除
	// 配置
	opt *SchedulerOpt
}

func (s *Scheduler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 跨域支持
	if s.opt.UseCORS {
		if s.opt.AllowOrigin == "" {
			w.Header().Add("Access-Control-Allow-Origin", "*")
		} else {
			w.Header().Add("Access-Control-Allow-Credentials", "true")
			w.Header().Add("Access-Control-Allow-Origin", s.opt.AllowOrigin)
		}
		w.Header().Add("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
	}
	// 返回数据格式
	w.Header().Add("Content-Type", "application/json")
	// 派遣函数
	if r.Method == "POST" && s.PostHandler != nil {
		s.PostHandler(w, r)
	}
	if r.Method == "GET" && s.GetHandler != nil {
		s.GetHandler(w, r)
	}
	if r.Method == "PUT" && s.PutHandler != nil {
		s.PutHandler(w, r)
	}
	if r.Method == "DELETE" && s.DeleteHandler != nil {
		s.DeleteHandler(w, r)
	}
	if r.Method == "PATCH" && s.PatchHandler != nil {
		s.PatchHandler(w, r)
	}
}
