package main

import (
	"chitchat/data"
	"net/http"
)

// Request结构包括以下部分：
// URL字段
// Header字段
// Body字段
// Form字段、PostForm字段、MultipartForm字段

// 为什么要以传值的方式将ResponseWriter传递给ServeHTTP？
// 答：接受Request结构指针的原因是为了让服务器能够察觉到处理器对Request结构的修改，其实ResponseWriter也是一样的，在查看net/http库
// 的源代码中会发现ResponseWriter实际上就是response这个非导出结构的接口，而ResponseWriter使用response结构时传递的也是指针。
// 所以，虽然ResponseWriter看上去像是一个值，但实际上是一个带有结构指针的接口。
// Responsewriter接口拥有以下3个方法：
// Write
// WriteHeader
// Header

// GET /err?msg=
// shows the error message page
func err(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	_, err := session(writer, request)
	if err != nil {
		generateHTML(writer, vals.Get("msg"), "layout", "public.navbar", "error")
	} else {
		generateHTML(writer, vals.Get("msg"), "layout", "private.navbar", "error")
	}
}

func index(writer http.ResponseWriter, request *http.Request) {
	threads, err := data.Threads()
	if err != nil {
		error_message(writer, request, "Cannot get threads")
	} else {
		_, err := session(writer, request)
		if err != nil {
			generateHTML(writer, threads, "layout", "public.navbar", "index")
		} else {
			generateHTML(writer, threads, "layout", "private.navbar", "index")
		}
	}
}
