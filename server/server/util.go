package server

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

func validPassword(s string) error {
next:
	for name, classes := range map[string][]*unicode.RangeTable{
		"Großbuchstabe": 	{unicode.Upper, unicode.Title},
		"Kleinbuchstabe": 	{unicode.Lower},
		"Zahl":    			{unicode.Number, unicode.Digit},
		//"Sonderzeichen":    {unicode.Space, unicode.Symbol, unicode.Punct, unicode.Mark},
	} {
		for _, r := range s {
			if unicode.IsOneOf(classes, r) {
				continue next
			}
		}
		return fmt.Errorf("Passwort muss 1 %s enthalten", name)
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}