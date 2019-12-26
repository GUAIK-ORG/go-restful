package restful

import "net/http"

type HttpHandler func(http.ResponseWriter, *http.Request)
type RestHandler func(http.ResponseWriter, *http.Request, map[string]interface{}, *Response)
