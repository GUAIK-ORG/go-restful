package restful

import "net/http"

type Filter interface {
	Processor(*http.Request, map[string]interface{}) (map[string]interface{}, error)
}
