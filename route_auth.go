package main

import (
	"chitchat/data"
	"net/http"
)

// GET /login
// Show the login page
func login(writer http.ResponseWriter, request *http.Request) {
	t := parseTemplateFiles("login.layout", "public.navbar", "login")
	t.Execute(writer, nil)
}

// GET /signup
// Show the signup page
func signup(writer http.ResponseWriter, request *http.Request) {
	generateHTML(writer, nil, "login.layout", "public.navbar", "signup")
}

// POST /signup
// Create the user account
func signupAccount(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		danger(err, "Cannot parse form")
	}
	user := data.User{
		Name:     request.PostFormValue("name"),
		Email:    request.PostFormValue("email"),
		Password: request.PostFormValue("password"),
	}
	if err := user.Create(); err != nil {
		danger(err, "Cannot create user")
	}
	http.Redirect(writer, request, "/login", 302)
}

// 当一个用户成功登录后，服务器必须在后续的请求中标示出这是一个已登录的用户
// 为了做到这一点，服务器会在响应的首部写入一个cookie，而客户端在接收这个cookie之后则会把它存储到浏览器里
// 下面的函数的作用就是对用户的身份进行验证，并在验证成功后像客户端返回一个cookie
// POST /authenticate
// Authenticate the user given the email and password
func authenticate(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	user, err := data.UserByEmail(request.PostFormValue("email"))
	if err != nil {
		danger(err, "Cannot find user")
	}
	if user.Password == data.Encrypt(request.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			danger(err, "Cannot create session")
		}
		// 没有为这个cookie设置过期时间，所以这是一个会话cookie，它将在浏览器关闭时自动被移除
		// HttpOnly字段被设置成true，意味着这个cookie只能通过HTTP或者HTTPS访问，
		// 无法通过JavaScript等非HTTP API进行访问
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(writer, &cookie)
		// 登录成功就跳转至主页
		http.Redirect(writer, request, "/", 302)
	} else {
		// TODO: 最好能提示密码或这用户名错误
		http.Redirect(writer, request, "/login", 302)
	}

}

// GET /logout
// Logs the user out
func logout(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("_cookie")
	if err != http.ErrNoCookie {
		warning(err, "Failed to get cookie")
		session := data.Session{Uuid: cookie.Value}
		session.DeleteByUUID()
	}
	http.Redirect(writer, request, "/", 302)
}
