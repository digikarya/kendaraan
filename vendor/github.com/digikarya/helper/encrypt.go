package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"github.com/google/uuid"
)

func GenerateUID() string{
	return uuid.New().String()
}
func CreateSHA(key string) string {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}



func CreateHashMd5(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encrypt(data string) ([]byte,error) {
	block, _ := aes.NewCipher([]byte(CreateHashMd5("PKfJ#6q$aS9")))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte(""),err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte(""),err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
	return ciphertext,nil
}

func Decrypt(data []byte) ([]byte,error) {
	key := []byte(CreateHashMd5("PKfJ#6q$aS9"))
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte(""),err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte(""),err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return []byte(""),err
	}
	return plaintext,nil
}


