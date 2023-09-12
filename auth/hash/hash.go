package hash

import (
	"golang.org/x/crypto/bcrypt"
)

func Make(s string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(s), 14)
	return string(bytes), err
}

func Check(unhashedString string, hashedString string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(unhashedString))
	return err == nil
}
