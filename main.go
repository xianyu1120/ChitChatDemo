package main

import (
	"net/http"
	"time"
)

func main() {
	p("ChitChat", version(), "started at", config.Address)
	//创建多路复用器
	mux := http.NewServeMux()
	//静态文件服务
	files := http.FileServer(http.Dir(config.Static))
	//收到、static/开头的url请求时，移除/static/，然后在public目录中查找被请求的文件
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	//index
	mux.HandleFunc("/", index)
	//err
	mux.HandleFunc("/err", err)

	//defind in route_auth.go
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/signup_account", signupAccount)
	mux.HandleFunc("/authenticate", authenticate)

	//defind in route_thread.go
	mux.HandleFunc("/thread/new", newThread)
	mux.HandleFunc("/thread/create", createThread)
	mux.HandleFunc("/thread/post", postThread)
	mux.HandleFunc("/thread/read", readThread)

	server := http.Server{
		Addr:              config.Address,
		Handler:           mux,
		ReadHeaderTimeout: time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:      time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes:    1 << 20,
	}
	server.ListenAndServe()
}
