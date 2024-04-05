package entity

import (
	"crypto/sha1"
	"os"
)

type User struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password []byte `json:"password"`
}

func (u *User) SetPassword(password string) {
	hash := sha1.New()
	hash.Write([]byte(password))
	u.Password = hash.Sum([]byte(os.Getenv("SALT"))) //TODO Salt
}
