package data

import "testing"

// test data
var users = []User{
	{
		Name:     "Peter Jones",
		Email:    "peter@gmail.com",
		Password: "peter_pass",
	},
	{
		Name:     "John Smith",
		Email:    "john@gmail.com",
		Password: "john_pass",
	},
}

func Test_setUp(t *testing.T) {
	//ThreadDeleteAll()
	SessionDeleteAll()
	UserDeleteAll()
}
