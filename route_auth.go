package main

import (
	"ChitChat/data"
	"net/http"
)

//GET /login
//登录
func login(w http.ResponseWriter, r *http.Request) {
	//解析模板文件
	t := parseTemplateFiles("login.layout", "public.navbar", "login")
	//执行
	t.Execute(w, nil)
}

//GET/logout
//退出
func logout(w http.ResponseWriter, r *http.Request) {
	//查看cookie,是否已登录
	cookie, err := r.Cookie("_cookie")
	if err != http.ErrNoCookie {
		//存在cookie,已登录
		warning(err, "Failed to get cookie")
		session := data.Session{Uuid: cookie.Value}
		session.DeleteByUUID()
	}
	http.Redirect(w, r, "/", 302)
}

//GET /signup
//注册界面
func signup(w http.ResponseWriter, r *http.Request) {
	generateHTML(w, nil, "login.layout", "public.navbar", "signup")
}

//post/signup
//注册账户
func signupAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() //解析表单，查看提交数据
	if err != nil {
		danger(err, "Cannot pars form")
	}
	user := data.User{
		Name:     r.PostFormValue("name"),
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}
	if err := user.Create(); err != nil {
		danger(err, "Cannot create user")
	}
	http.Redirect(w, r, "/login", 302) //重定向
}

//post /authenticate
//验证给定电子邮件和密码的用户,对用户身份进行认证，认证成功后返回一个cookie
func authenticate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析请求
	user, err := data.UserByEmail(r.PostFormValue("email"))
	if err != nil {
		danger(err, "Cannot find user")
	}
	if user.Password == data.Encrypt(r.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			danger(err, "Cannot create session")
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true, //只能通过HTTP或者https访问
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", 302)
	} else {
		http.Redirect(w, r, "/login", 302)
	}
}
