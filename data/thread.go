package data

import (
	"log"
	"time"
)

//Thread 论坛中的帖子
type Thread struct {
	Id        int
	Uuid      string
	Topic     string
	UserId    int
	CreatedAt time.Time
}

//Post 用户在帖子中的回复
type Post struct {
	Id        int
	Uuid      string
	Body      string
	UserId    int
	ThreadId  int
	CreatedAt time.Time
}

// format the CreatedAt date to display nicely on the screen
func (t *Thread) CreatedAtDate() string {
	return t.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

func (p *Post) CreatedAtDate() string {
	return p.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

//获取贴子数
func (t *Thread) NumReplies() (count int) {
	rows, err := DB.Query("select count(*)from posts where thread_id=?", t.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return
		}
	}
	_ = rows.Close()
	return
}
func (t *Thread) Posts() (posts []Post, err error) {
	rows, err := DB.Query("SELECT id, uuid, body, user_id, thread_id, created_at FROM posts where thread_id = ?", t.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}
		if err = rows.Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt); err != nil {
			return
		}
		posts = append(posts, post)
	}
	_ = rows.Close()
	return
}
func (u *User) CreateThread(topic string) (conv Thread, err error) {
	statement := "insert into threads(uuid,topic,user_id,created_at)values(?,?,?,?)"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(createUUID(), topic, u.Id, time.Now()).Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt)
	return
}
func (u *User) CreatePost(conv Thread, body string) (post Post, err error) {
	statement := "insert into posts (uuid, body, user_id, thread_id, created_at) values (?, ?,?, ?, ?)"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(createUUID(), body, u.Id, conv.Id, time.Now()).Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt)
	return
}
func Threads() (threads []Thread, err error) {
	rows, err := DB.Query("select id,uuid,topic,user_id,created_at from threads order by created_at desc ")
	if err != nil {
		log.Print("query failed")
		return
	}
	for rows.Next() {
		th := Thread{}
		if err = rows.Scan(&th.Id, &th.Uuid, &th.Topic, &th.UserId, &th.CreatedAt); err != nil {
			return
		}
		threads = append(threads, th)
	}
	defer rows.Close()
	return
}
func ThreadByUUID(uuid string) (conv Thread, err error) {
	conv = Thread{}
	err = DB.QueryRow("SELECT id, uuid, topic, user_id, created_at FROM threads WHERE uuid = ?", uuid).
		Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt)
	return
}

// Get the user who started this thread
func (t *Thread) User() (user User) {
	user = User{}
	_ = DB.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = ?", t.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

// Get the user who wrote the post
func (p *Post) User() (user User) {
	user = User{}
	_ = DB.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = ?", p.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}
