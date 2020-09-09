package main

import (
	"ChitChat/data"
	"log"
	"net/http"
)

func err(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	_, err := session(w, r)
	if err != nil {
		generateHTML(w, vals.Get("msg"), "layout", "public.navbar", "error")
	} else {
		generateHTML(w, vals.Get("msg"), "layout", "private.navbar", "error")
	}
}

//index 负责生成html，并将其写入到responseWriter中
func index(w http.ResponseWriter, r *http.Request) {
	//显示所有已有帖子
	//判断用户是否登录
	threads, err := data.Threads()
	if err != nil {
		log.Print("failed")
		error_message(w, r, "Cannot get threads")
	} else {
		_, err := session(w, r)
		if err != nil {
			generateHTML(w, r, "layout", "public.navbar", "index")
		} else {
			generateHTML(w, threads, "layout", "private.navbar", "index")
		}
	}
}
