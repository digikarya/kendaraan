package helper

import (
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	Password string
	encode string
}

func (r *Password) init(password string){
	r.Password = password
	r.encode = CreateHashMd5(password+"QVP7sQLtLeRG4c5N")
}

func  (r *Password)  HashPassword(pass string) (string, error) {
	r.init(pass)
	bytes, err := bcrypt.GenerateFromPassword([]byte(r.encode), 14)
	return string(bytes), err
}

func (r *Password) Check(password, hash string) bool {
	r.init(password)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(r.encode))
	return err == nil
}