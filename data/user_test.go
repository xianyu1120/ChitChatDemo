package data

import (
	"log"
	"testing"
)

func TestUser_CreateSession(t *testing.T) {
	if session, err := users[0].CreateSession(); err != nil {
		panic(err)
	} else {
		log.Printf("%v\n", session)
	}
}
func TestUser_Session(t *testing.T) {
	user := User{
		Id: 0,
	}
	session, err := user.Session()
	if err != nil {
		panic(err)
	} else {
		log.Printf("%v\n", session)
	}
	if valid, err := session.Check(); valid != true || err != nil {
		panic(err)
	}

	//TestDeleteByUUID
	if err = session.DeleteByUUID(); err != nil {
		panic(err)
	} else {
		valid, err := session.Check()
		if valid == true || err == nil {
			panic(err)
		}
	}
}
func TestSession_User(t *testing.T) {
	s := Session{UserId: 0}
	user, err := s.User()
	if err != nil {
		panic(err)
	}
	log.Print(user)

}
func TestSessionDeleteAll(t *testing.T) {
	err := SessionDeleteAll()
	if err != nil {
		panic(err)
	}
}
func TestUser_Create(t *testing.T) {
	if err := users[0].Create(); err != nil {
		t.Error(err, "cannot create user")
	}
	u, err := UserByEmail(users[0].Email)
	if err != nil {
		t.Error(err, "User not created.")
	}
	if users[0].Email != u.Email {
		t.Errorf("User retrieved is not the same as the one created.")
	}
	if err = users[0].Delete(); err != nil {
		panic(err)
	}
	_, err = UserByEmail(users[0].Email)
	if err != nil {
		t.Error("user is delete")
	}
}
