package main

import (
	"ChitChat/data"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Configuration struct {
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
}

var config Configuration
var logger *log.Logger

//标准输出
func p(a ...interface{}) {
	fmt.Println(a)
}

//初始化
func init() {
	loadConfig() //解析配置
	//创建错误日志
	file, err := os.OpenFile("chitchat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("open log file failed", err)
	}
	logger = log.New(file, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
}

//服务器配置解析
func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
}

//重定向到错误界面
func error_message(w http.ResponseWriter, r *http.Request, msg string) {
	url := []string{"/err?msg=", msg}
	http.Redirect(w, r, strings.Join(url, ""), 302)

}
func session(w http.ResponseWriter, r *http.Request) (sess data.Session, err error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		sess = data.Session{Uuid: cookie.Value}
		if ok, _ := sess.Check(); !ok {
			err = errors.New("Invalid session")
		}
	}
	return
}

//解析template 模板
func parseTemplateFiles(filenames ...string) (t *template.Template) {
	var files []string
	t = template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	t = template.Must(t.ParseFiles(files...))
	return
}

//生成html
func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(w, "layout", data)
}

//日志初始化
func info(args ...interface{}) {
	logger.SetPrefix("INFO")
	logger.Println(args...)
}

//危险
func danger(args ...interface{}) {
	logger.SetPrefix("ERROR") //日志输出前缀“error"
	logger.Println(args...)
}

//错误
func warning(args ...interface{}) {
	logger.SetPrefix("WARNING")
	logger.Println(args...)
}
func version() string {
	return "0.1"
}
