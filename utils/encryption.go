package utils

import "golang.org/x/crypto/bcrypt"

func Encrypt(password string) string {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashedPass)
}

func ComparePass(pass1, pass2 string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pass1), []byte(pass2))
	if err != nil {
		return false
	} else {
		return true
	}

}
