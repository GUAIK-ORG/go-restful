package cmd

import (
	filter "go-restful/filters"
	"go-restful/parser"
	"go-restful/pkg/restful"
	"net/http"
)

// 设置Session(登陆处理)
var CreateSessionHandler = restful.NewHandler(
	// 处理函数
	func(w http.ResponseWriter, r *http.Request, params map[string]interface{}, resp *restful.Response) {
		// 模拟数据库验证操作
		uid := "10001"
		token := "xxxxxxxx"
		if params["email"].(string) == "demo@guaik.org" && params["passwd"] == "hello!" {
			// 设置cookie
			http.SetCookie(w, &http.Cookie{Name: "uid", Value: uid})
			http.SetCookie(w, &http.Cookie{Name: "token", Value: token})
			// 返回uid和token
			resp.Success(map[string]interface{}{
				"uid":   uid,
				"token": token,
			})
		} else {
			// 使用已配置的错误信息
			resp.UseError("SESSION.10000")
		}
	},
	// 接口配置
	restful.HandlerOpts{
		// 配置接口错误信息
		MakeErrorFunc: func(err *restful.Errors) {
			err.NewError("SESSION.10000", "email or passwd error")
			err.Translate("SESSION.10000", "cn", "邮箱或密码错误") // 中文翻译
		},
		// 设置解析器
		ParseFunc: parser.JsonParser,
		// 配置过滤器
		Filters: []restful.Filter{
			&filter.CheckParams{
				// 参数检查
				Params: map[string]interface{}{
					// 正则校验
					"email": filter.FieldRegexp(`^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$`),
					// 6～12长度字符串校验
					"passwd": filter.FieldString().SetLength(6, 12),
				},
			},
		},
	},
)

// 删除Session(注销)
var DeleteSessionHandler = restful.NewHandler(
	// 处理函数
	func(w http.ResponseWriter, r *http.Request, params map[string]interface{}, resp *restful.Response) {
		_, err := r.Cookie("uid")
		if err != nil {
			resp.UseError("SESSION.10001")
		} else {
			resp.Success(nil)
		}
	},
	// 接口配置
	restful.HandlerOpts{
		// 配置接口错误信息
		MakeErrorFunc: func(err *restful.Errors) {
			err.NewError("SESSION.10001", "delete session error")
			err.Translate("SESSION.10001", "cn", "删除会话失败")
		},
		Filters: []restful.Filter{
			// 该接口需要验证token，如果token无效将不被执行
			&filter.CheckToken{},
		},
	},
)
