package data

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/chitchat?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	return
}

//创建随机uuid
func createUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("cannot generate UUID", err)
	}
	//0x40是RFC 4122的保留变体
	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

//加密给定字符串
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}
