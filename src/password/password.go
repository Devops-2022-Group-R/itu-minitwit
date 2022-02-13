// A small package for generating password hashes and checking if a plaintext
// password matches a given password hash. Based on werkzeug from Python with Go
// code modified from saggit/main.go on GitHub Gist: https://gist.github.com/saggit/19c4404e9a20d54fdcf1
package password

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	keyLength         = 32
	saltLength        = 16
	saltChars         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	defaultIterations = 260000
)

func generateSalt() string {
	var bytes = make([]byte, saltLength)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = saltChars[v%byte(len(saltChars))]
	}
	return string(bytes)
}

func hashInternal(salt string, password string, iterations int) string {
	hash := pbkdf2.Key([]byte(password), []byte(salt), iterations, keyLength, sha256.New)
	return hex.EncodeToString(hash)
}

func GeneratePasswordHash(password string) string {
	salt := generateSalt()
	hash := hashInternal(salt, password, defaultIterations)
	return fmt.Sprintf("pbkdf2:sha256:%d$%s$%s", defaultIterations, salt, hash)
}

func CheckPasswordHash(password string, hash string) bool {
	if strings.Count(hash, "$") < 2 {
		return false
	}
	methodSaltHash := strings.Split(hash, "$")
	method, salt, hash := methodSaltHash[0], methodSaltHash[1], methodSaltHash[2]
	iterations, _ := strconv.Atoi(strings.Split(method, ":")[2])
	return hash == hashInternal(salt, password, iterations)
}
