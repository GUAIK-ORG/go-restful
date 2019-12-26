package restful

import (
	"fmt"
	"net/http"
)

type Restful struct {
	router map[string]*Scheduler
	mux    *http.ServeMux
	defOpt *SchedulerOpt
}

// 创建Restful服务
func NewRestful() *Restful {
	return &Restful{
		router: make(map[string]*Scheduler),
		defOpt: nil,
		mux:    http.NewServeMux(),
	}
}

func (rf *Restful) SetDefOpt(opt *SchedulerOpt) {
	rf.defOpt = opt
}

func (rf *Restful) Scheduler(path string, opt *SchedulerOpt) {
	if _, ok := rf.router[path]; !ok {
		if opt == nil && rf.defOpt != nil {
			opt = rf.defOpt
		}
		rf.router[path] = &Scheduler{opt: opt}
		rf.mux.Handle(path, rf.router[path])
	}
}

func (rf *Restful) Post(path string, handler HttpHandler) {
	rf.Scheduler(path, nil)
	if rf.router[path].PostHandler == nil {
		rf.router[path].PostHandler = handler
	}

}

func (rf *Restful) Get(path string, handler HttpHandler) {
	rf.Scheduler(path, nil)
	if rf.router[path].GetHandler == nil {
		rf.router[path].GetHandler = handler
	}

}

func (rf *Restful) Put(path string, handler HttpHandler) {
	rf.Scheduler(path, nil)
	if rf.router[path].PutHandler == nil {
		rf.router[path].PutHandler = handler
	}

}

func (rf *Restful) Delete(path string, handler HttpHandler) {
	rf.Scheduler(path, nil)
	if rf.router[path].DeleteHandler == nil {
		rf.router[path].DeleteHandler = handler
	}
}

func (rf *Restful) Start(port uint32) {
	http.ListenAndServe(fmt.Sprintf(":%d", port), rf.mux)
}
