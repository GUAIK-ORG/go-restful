package filter

import (
	"errors"
	"net/http"
	"reflect"
	"regexp"
	"unicode/utf8"

	"github.com/golang/glog"
)

type CkField struct {
	T      string // 类型
	R      string // 正则表达式
	Sparse bool   // 为空则不处理
	// 长度检查
	MaxLen int
	MinLen int
}

func (c *CkField) UseSparse(b bool) *CkField {
	c.Sparse = b
	return c
}

func (c *CkField) SetLength(min, max int) *CkField {
	c.MinLen = min
	c.MaxLen = max
	return c
}

func FieldString() *CkField {
	return &CkField{T: "string", MinLen: -1, MaxLen: -1, Sparse: false}
}

func FieldInt() *CkField {
	return &CkField{T: "int", MinLen: -1, MaxLen: -1, Sparse: false}
}

func FieldFloat64() *CkField {
	return &CkField{T: "float64", MinLen: -1, MaxLen: -1, Sparse: false}
}

func FieldBool() *CkField {
	return &CkField{T: "bool", MinLen: -1, MaxLen: -1, Sparse: false}
}

func FieldRegexp(reg string) *CkField {
	return &CkField{R: reg, MinLen: -1, MaxLen: -1, Sparse: false}
}

type CkList struct {
	MinLimit   int32         // 最小限制（使用模版时生效） -1：不处理
	MaxLimit   int32         // 最大限制（使用模版时生效） -1：不处理
	isTemplete bool          // 启用后，L中的第一个项会被当成列表其他元素的模板
	L          []interface{} // 列表内容
}

func (c *CkList) UseTemplete(min, max int32) *CkList {
	c.MinLimit = min
	c.MaxLimit = max
	c.isTemplete = true
	return c
}

func List(l []interface{}) *CkList {
	return &CkList{MinLimit: -1, MaxLimit: -1, L: l, isTemplete: false}
}

type CheckParams struct {
	Params map[string]interface{}
}

// 字典检查
func (c *CheckParams) checkFunc(one interface{}, two interface{}) bool {
	if one == nil || two == nil {
		glog.Error("filter@check faild == nil")
		return false
	}
	if t2, ok := two.(map[string]interface{}); ok {
		// 判断是否是字典
		if t1, ok := one.(map[string]interface{}); ok {
			for k, v := range t2 { // 遍历字典
				if _, ok := t1[k]; !ok {
					// 参数缺少字段
					// 判断是否是稀疏字段
					if t, ok := v.(*CkField); ok && t.Sparse {
						continue
					}
					return false
				}
				if !c.checkFunc(t1[k], v) {
					return false
				}
			}
		} else {
			glog.Error("filter@check map error")
			return false // 类型不匹配
		}
	} else if t2, ok := two.(*CkList); ok {
		if t1, ok := one.([]interface{}); ok {
			if t2.isTemplete {
				// 使用模版
				if len(t2.L) != 1 {
					glog.Error("filter@check use templete but len(L) != 1")
					return false
				}
				// 长度检查
				if t2.MinLimit >= 0 && int32(len(t1)) < t2.MinLimit {
					glog.Error("filter@check min limit error")
					return false
				}
				if t2.MaxLimit >= 0 && int32(len(t1)) > t2.MaxLimit {
					glog.Error("filter@check max limit error")
					return false
				}
				for _, v := range t1 {
					if !c.checkFunc(v, t2.L[0]) {
						return false
					}
				}
			} else {
				// 不使用模版(精确匹配内容)
				if len(t2.L) != len(t1) {
					glog.Error("filter@check list length error")
					return false
				}
				for i, v := range t2.L {
					if !c.checkFunc(t1[i], v) {
						return false
					}
				}
			}
		} else {
			return false // 类型不匹配
		}
	} else if t2, ok := two.(*CkField); ok {
		// 检查字段
		if t2.R != "" {
			// 通过正则校验字段
			if str, ok := one.(string); ok {
				if ok, _ := regexp.MatchString(t2.R, str); !ok {
					glog.Error("filter@check regexp error : ", str)
					return false
				}
			} else {
				// 使用正则表达式的字段必须为string类型，否则返回false
				glog.Error("filter@check regexp type not string")
				return false
			}
		} else if t2.T != reflect.TypeOf(one).String() {
			glog.Error("filter@check type error")
			return false
		} else {
			// 字符串长度检查
			if str, ok := one.(string); ok {
				len := utf8.RuneCountInString(str)
				// glog.Info("check string length: min:", t2.MinLen, " max:", t2.MaxLen, " ->", len)
				if t2.MinLen != -1 && len < t2.MinLen {
					glog.Error("filter@check min len error")
					return false
				}
				if t2.MaxLen != -1 && len > t2.MaxLen {
					glog.Error("filter@check max len error")
					return false
				}
			}
		}
	} else {
		// 非map,list,CKField的类型不支持直接返回false
		glog.Error("filter@check other error")
		return false
	}
	return true
}

func (c *CheckParams) Processor(r *http.Request, in map[string]interface{}) (out map[string]interface{}, err error) {
	out = in
	if !c.checkFunc(in, c.Params) {
		err = errors.New("filter: params key not exist or type error")
	}
	return
}
