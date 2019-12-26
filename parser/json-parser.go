package parser

import (
	"encoding/json"
	"errors"
)

func JsonParser(urlParams map[string][]string, data []byte) (params map[string]interface{}, err error) {
	params = make(map[string]interface{})
	if len(data) != 0 {
		// 处理POST和PUT请求的Body(优先处理)
		// 如果GET和DELETE请求也使用Body，
		// 那么自动放弃处理URL参数
		err = json.Unmarshal(data, &params)
		if err != nil {
			return
		}
	} else if urlParams != nil && len(urlParams) != 0 {
		// 处理GET和DELETE请求
		for k, v := range urlParams {
			params[k] = v[0]
		}
	} else {
		err = errors.New("Parser: body empty")
	}
	return
}
