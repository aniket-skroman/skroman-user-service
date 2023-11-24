package utils

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func RandomString(lenght int) (ans string) {
	chars := "abcdefghijklmnopqrstuvwxyz"
	for len(ans) < lenght {
		ans += string(chars[rand.Intn(len(chars))])
	}

	return
}

func RandomInt() (ans string) {
	t := fmt.Sprint(time.Now().Nanosecond())
	ans = t[:7]
	return
}

func Hash_password(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Fatal(err)
	}

	return string(bytes)
}

func Compare_password(hashedPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)

	if err != nil {
		log.Println(err)
		return false
	}
	return true

}
