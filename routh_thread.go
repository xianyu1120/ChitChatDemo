package main

import (
	"ChitChat/data"
	"fmt"
	"net/http"
)

//GET /threads/new
//显示新的帖子
func newThread(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		generateHTML(w, nil, "layout", "private.navbar", "new.thread")
	}
}

//Post /signup
//创建新的帖子
func createThread(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		err = r.ParseForm()
		if err != nil {
			danger(err, "Cannot parse form")
		}
		user, err := sess.User()
		if err != nil {
			danger(err, "Cannot get user from session")
		}
		topic := r.PostFormValue("topic")
		if _, err := user.CreateThread(topic); err != nil {
			danger(err, "Cannot create thread")
		}
		http.Redirect(w, r, "/", 302)
	}
}

//GET /thread/read
//显示帖子的详细信息，包括帖子和撰写帖子的表单
func readThread(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	uuid := vals.Get("id")
	thread, err := data.ThreadByUUID(uuid)
	if err != nil {
		error_message(w, r, "Cannot read thread")
	} else {
		_, err := session(w, r)
		if err != nil {
			generateHTML(w, &thread, "layout", "public.navbar", "public.thread")
		} else {
			generateHTML(w, &thread, "layout", "private.navbar", "private.thread")
		}
	}
}

//post /thread/post
//创建信息
func postThread(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		err = r.ParseForm()
		if err != nil {
			danger(err, "Cannot parse form")
		}
		user, err := sess.User()
		if err != nil {
			danger(err, "Cannot get user from session")
		}
		body := r.PostFormValue("body")
		uuid := r.PostFormValue("uuid")
		thread, err := data.ThreadByUUID(uuid)
		if err != nil {
			error_message(w, r, "cannot read thread")
		}
		if _, err := user.CreatePost(thread, body); err != nil {
			danger(err, "Cannot create post")
		}
		url := fmt.Sprint("/thread/read?id=", uuid)
		http.Redirect(w, r, url, 302)
	}
}
