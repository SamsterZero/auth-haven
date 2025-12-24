package utils

import "golang.org/x/crypto/bcrypt"

// Hash hashes a plain password and returns the hash or an error
func Hash(pswd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pswd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword compares a plain password with its hashed value
func CheckPassword(hash, pswd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pswd))
}
