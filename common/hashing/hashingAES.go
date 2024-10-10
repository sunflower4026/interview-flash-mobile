package hashing

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func Decrypt(ciphertext string, key, iv []byte) (result []byte, err error) {
	//Remove base64 encoding:
	cipherText, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return
	}

	// Create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	// IF the length of the cipherText is less than 16 Bytes:
	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short")
		return
	}

	// Decrypt the message
	stream := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(cipherText))
	stream.CryptBlocks(decrypted, cipherText)

	result = pkcs5Trimming(decrypted)

	return
}

func Encrypt(plaintext []byte, key, iv []byte) (encoded string, err error) {
	bPlaintext := pkcs5Padding(plaintext, aes.BlockSize, len(string(plaintext)))

	c, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	cipherText := make([]byte, len(bPlaintext))

	stream := cipher.NewCBCEncrypter(c, iv)
	stream.CryptBlocks(cipherText, bPlaintext)
	return base64.StdEncoding.EncodeToString(cipherText), err
}

func pkcs5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func Newencrypt(message []byte, key []byte) (encoded string, err error) {
	//Create byte array from the input string
	plainText := message

	//Create a new AES cipher using the key
	block, err := aes.NewCipher(key)

	//IF NewCipher failed, exit:
	if err != nil {
		return
	}

	//Make the cipher text a byte array of size BlockSize + the length of the message
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	//iv is the ciphertext up to the blocksize (16)
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	//Encrypt the data:
	stream := cipher.NewCBCEncrypter(block, iv)
	stream.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	//Return string encoded in base64
	return base64.RawStdEncoding.EncodeToString(cipherText), err
}

func Newdecrypt(secure string, key []byte) (decoded []byte, err error) {
	//Remove base64 encoding:
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)

	//IF DecodeString failed, exit:
	if err != nil {
		return
	}

	//Create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher(key)

	//IF NewCipher failed, exit:
	if err != nil {
		return
	}

	//IF the length of the cipherText is less than 16 Bytes:
	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short")
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	//Decrypt the message
	stream := cipher.NewCBCDecrypter(block, iv)
	stream.CryptBlocks(cipherText, cipherText)

	return cipherText, err
}
