package data

import (
	"log"
	"time"
)

type User struct {
	Id        int
	Uuid      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

//Session 用户当前登录会话
type Session struct {
	Id        int
	Uuid      string //存储一个随机生成的唯一id
	Email     string
	UserId    int //用户表中存储用户信息的行的id
	CreatedAt time.Time
}

//为现有用户创建一个新的会话
func (u *User) CreateSession() (session Session, err error) {
	statement := "insert into sessions (uuid, email, user_id, created_at) values (?, ?, ?, ?)"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		log.Printf("err:%v\n", err)
		return
	}
	defer stmt.Close()
	uuid := createUUID()
	email := u.Email
	uid := u.Id
	createAt := time.Now()
	result, err := stmt.Exec(uuid, email, uid, createAt) //执行插入
	id, _ := result.LastInsertId()                       //新插入数据id

	session = Session{
		Id:        int(id),
		Uuid:      uuid,
		Email:     email,
		UserId:    uid,
		CreatedAt: createAt,
	}
	return
}

//获取现有用户的会话
func (u *User) Session() (session Session, err error) {
	session = Session{}
	err = DB.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE user_id = ?", u.Id).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
}

//检查会话在数据库中是否有效
func (s *Session) Check() (valid bool, err error) {
	err = DB.QueryRow("select id,uuid,email,user_id,created_at from sessions where uuid=?", s.Uuid).
		Scan(&s.Id, &s.Uuid, &s.Email, &s.UserId, &s.CreatedAt)
	if err != nil {
		valid = false //无效
		return
	}
	if s.Id != 0 {
		valid = true
	}
	return
}

//将会话从数据库中删除
func (s *Session) DeleteByUUID() (err error) {
	statement := "delete from sessions where uuid=?"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(s.Uuid)
	return
}

//从会话中得到用户
func (s *Session) User() (user User, err error) {
	user = User{}
	err = DB.QueryRow("select id,uuid,name,email,created_at from users where id=?", s.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

//从数据库中删除所有会话
func SessionDeleteAll() (err error) {
	statement := "delete  from sessions"
	_, err = DB.Exec(statement)
	return
}

//创建一个新的用户，保存到数据库中
func (u *User) Create() (err error) {
	statement := "insert into users(uuid,name,email,password,created_at)values(?,?,?,?,?)"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		log.Print("创建用户失败")
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(createUUID(), u.Name, u.Email, Encrypt(u.Password), time.Now()) //执行
	if err != nil {
		log.Print("插入用户数据失败")
	}
	return
}

//删除一个用户
func (u *User) Delete() (err error) {
	statement := "delete from users where id=?"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(u.Id)
	return
}

//更新用户在数据库中的信息
func (u *User) Update() (err error) {
	statement := "update users set name=?,email=? where id=?"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(u.Id, u.Name, u.Email)
	return
}

//删除所有用户
func UserDeleteAll() (err error) {
	statement := "delete from users"
	_, err = DB.Exec(statement)
	return
}

//获取所有用户
func Users() (users []User, err error) {
	rows, err := DB.Query("select id,uuid,name,email,password,created_at from users")
	if err != nil {
		return
	}
	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt); err != nil {
			return
		}
		users = append(users, user)
	}
	rows.Close()
	return
}

//通过Email得到一个用户
func UserByEmail(email string) (user User, err error) {
	user = User{}
	err = DB.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE email = ?", email).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

//根据uuid获取一个用户
func UserByUUID(uuid string) (user User, err error) {
	user = User{}
	err = DB.QueryRow("select id,uuid,name,email,password,created_at from users where uuid=?", uuid).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}
