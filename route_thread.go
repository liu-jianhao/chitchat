package main

import (
	"chitchat/data"
	"fmt"
	"net/http"
)

// 发帖子的处理函数
// GET /threads/new
// Show the new thread form page
func newThread(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err != nil {
		// 要想发帖子，必须先登录，即跳转到登录页面
		http.Redirect(writer, request, "/login", 302)
	} else {
		// 如果已经登录了，就生成发帖子的页面
		generateHTML(writer, nil, "layout", "private.navbar", "new.thread")
	}
}

// 创建新帖子的处理函数
// POST /signup
// Create the user account
func createThread(writer http.ResponseWriter, request *http.Request) {
	// 检查用户是否登录，并获得会话信息
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {
		err = request.ParseForm()
		if err != nil {
			danger(err, "Cannot parse form")
		}
		// 从sess中获得用户结构
		user, err := sess.User()
		if err != nil {
			danger(err, "Cannot get user from session")
		}
		topic := request.PostFormValue("topic")
		// 调用用户中的方法CreateThread，将帖子写入数据库
		if _, err := user.CreateThread(topic); err != nil {
			danger(err, "Cannot create thread")
		}
		// 创建成功后回到主页面
		http.Redirect(writer, request, "/", 302)
	}
}

// 读取帖子的处理函数
// GET /thread/read
// Show the details of the thread, including the posts and the form to write a post
func readThread(writer http.ResponseWriter, request *http.Request) {
	// 解析URL中请求的帖子
	vals := request.URL.Query()
	uuid := vals.Get("id")
	// 根据uuid从数据库中获取帖子的信息
	thread, err := data.ThreadByUUID(uuid)
	if err != nil {
		error_message(writer, request, "Cannot read thread")
	} else {
		// 根据是否登录分为可以回帖和不可回帖两种页面
		_, err := session(writer, request)
		if err != nil {
			generateHTML(writer, &thread, "layout", "public.navbar", "public.thread")
		} else {
			generateHTML(writer, &thread, "layout", "private.navbar", "private.thread")
		}
	}
}

// 回帖的处理函数
// POST /thread/post
// Create the post
func postThread(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {
		err = request.ParseForm()
		if err != nil {
			danger(err, "Cannot parse form")
		}
		user, err := sess.User()
		if err != nil {
			danger(err, "Cannot get user from session")
		}
		body := request.PostFormValue("body")
		uuid := request.PostFormValue("uuid")
		thread, err := data.ThreadByUUID(uuid)
		if err != nil {
			error_message(writer, request, "Cannot read thread")
		}
		if _, err := user.CreatePost(thread, body); err != nil {
			danger(err, "Cannot create post")
		}
		url := fmt.Sprint("/thread/read?id=", uuid)
		http.Redirect(writer, request, url, 302)
	}
}
