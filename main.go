package main

import (
	"flag"
	"go-restful/cmd"
	"go-restful/pkg/restful"
	"runtime"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer func() {
		glog.Flush()
	}()
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Restful初始化
	server := restful.NewRestful()
	// 设置默认接口配置信息
	server.SetDefOpt(&restful.SchedulerOpt{UseCORS: true, AllowOrigin: "*"})

	/* Session */
	// 添加用户session（登陆操作）
	server.Post("/session", cmd.CreateSessionHandler)
	// 删除用户session（注销操作）
	server.Delete("/session", cmd.DeleteSessionHandler)

	// 运行Restful接口服务
	server.Start(8080)
}
