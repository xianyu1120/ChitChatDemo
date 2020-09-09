package data

import (
	"log"
	"testing"
)

func Test_CreateThread(t *testing.T) {
	err := DB.Ping()
	if err != nil {
		log.Print("Cannot　connect to the databases")
	}
	if err := users[0].Create(); err != nil {
		t.Error(err, "Cannot create user.")
	}
	conv, err := users[0].CreateThread("My first thread")
	if err != nil {
		t.Error(err, "Cannot create thread")
	}
	if conv.UserId != users[0].Id {
		t.Error("User not linked with thread")
	}
}
func Test_Thread(t *testing.T) {
	threads, err := Threads()
	if err != nil {
		panic(err)
	}
	log.Printf("%v\n", threads)
}
