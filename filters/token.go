package filter

import (
	"net/http"

	"github.com/golang/glog"
)

type CheckToken struct {
}

func (c *CheckToken) Processor(r *http.Request, in map[string]interface{}) (out map[string]interface{}, err error) {
	var (
		uid, token *http.Cookie
	)
	out = in
	uid, err = r.Cookie("uid")
	if err != nil {
		return
	}
	token, err = r.Cookie("token")
	if err != nil {
		return
	}
	// 这里添加验证token的代码，如果err != nil
	// 则filter验证失败
	// ......
	glog.Info("uid: %v , token: %v", uid.Value, token.Value)
	return
}
