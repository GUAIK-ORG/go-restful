# 🚀GO-Restful框架

## 快速开始

### 安装

`git clone https://github.com/GUAIK-ORG/go-restful.git`

### 运行

`go run main.go -log_dir=log -alsologtostderr`

### 测试

`./test/session.html`提供了一个js的登陆测试用例，请双击运行。测试用的邮箱和密码为：`email:demo@guaik.org passwd:hello!`

### 负载测试

新建文件：abtest.txt，内容如下:

```json
{
 "email": "demo@guaik.org",
 "passwd": "hello!"
}
```

ab测试指令：`ab -n 10000 -c 100 -p "./abtest.txt" -T "application/json" -H "Content-Type: application/json" http://localhost:8080/session`

本机测试：（有条件的可用服务器测试）

| 参数 | 配置 |
|---|---|
| OS | MacBook Pro (13-inch, Late 2016, Four Thunderbolt 3 Ports)|
| CPU | 2.9 GHz 双核Intel Core i5 |
| RAM | 8 GB 2133 MHz LPDDR3 |

测试结果：

![ab-load](https://gitee.com/GuaikOrg/go-restful/raw/master/docs/ab-load.png)

---

## 框架介绍

![arch](https://gitee.com/GuaikOrg/go-restful/raw/master/docs/arch.png)

框架代码在`pkg/restful`目录下

go-restful标准化了Restful接口开发，提供了`post delete put get`四种操作方式。

在`./cmd`目录下`session.go`实现了一个标准的Restful处理者，可参考使用。

框架提供了标准的返回数据：当status为0时代表操作成功，并且可在body中获取返回数据。

在handler中设置成功状态：

```go
resp.Success(map[string]interface{} {
    "uid":   uid,
    "token": token,
})
```

客户端接收到的数据为：

```json
{"status": 0, "error_code": 0, "error_msg": null, "body": {"token": "xxxxxxxx", "uid": "10001"}}
```

框架提供了多语言的错误信息，可通过配置的形式注册错误信息：

```go
restful.HandlerOpts{
    // 配置接口错误信息
    MakeErrorFunc: func(err *restful.Errors){
        err.NewError(1000, "email or passwd error")
        err.Translate(1000, "cn", "邮箱或密码错误") // 中文翻译
    },
},
```

客户端接收到的数据为：

```json
{"status": -1, "error_code": 1000, "error_msg": {"cn": "邮箱或密码错误", "en": "email or passwd error"}, "body": null}
```

框架可自定义请求解析器，默认提供了json格式解析在`./parser/json-parser.go`中。

框架支持过滤器队列，对请求数据进行预处理，在目录`./filters`目录下默认提供了两个过滤器。

check.go : 负责参数格式校验，支持string，float64，int64，bool，[]interface{}，正则表达式校验。

token.go : 用来校验访问令牌信息。（需结合缓存和数据库进行修改）。

将过滤器用于处理者：只要有任何一个过滤器`error != nil`，之后的过滤器将不会被执行，请求将被丢弃。

```go
restful.HandlerOpts{
    Filters: []restful.Filter{
        // 1、该接口需要验证token，如果token无效将不被执行
        &filter.CheckToken{},
        // 2、校验参数
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
```
