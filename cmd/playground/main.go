package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	b, _ := os.ReadFile("source_pass.txt")

	s := string(b)
	passwords := strings.Split(s, "\n")

	for i, p := range passwords {
		passwords[i] = fmt.Sprintf("%s: %s", p, HashPassword(p))
	}

	os.WriteFile("bcrypted_pass.txt", []byte(strings.Join(passwords, "\n")), 0644)
}

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(hash)
}
