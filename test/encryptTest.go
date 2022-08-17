package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	pwd := "123"
	encryptedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	enPwd := string(encryptedPwd)
	fmt.Println(enPwd)
	if bcrypt.CompareHashAndPassword([]byte(enPwd), []byte(pwd)) != nil {
		fmt.Println("false")
	} else {
		fmt.Println("true")
	}
}
