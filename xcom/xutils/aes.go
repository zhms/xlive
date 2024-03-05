package xutils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func padPKCS7(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func unpadPKCS7(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:length-unpadding]
}

func AesEncrypt(plaintext []byte, key []byte) ([]byte, error) {
	plaintext = padPKCS7(plaintext, aes.BlockSize)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	base64String := base64.StdEncoding.EncodeToString(ciphertext)
	return []byte(base64String), nil
}

func AesDecrypt(ciphertext []byte, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(string(ciphertext))
	ciphertext = data
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = unpadPKCS7(ciphertext)
	return ciphertext, nil
}

/*
	key := []byte("examplekey123456") // 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256
	plaintext := []byte{72, 101, 108, 108, 111, 44, 32, 87, 111, 114, 108, 100, 33}
	fmt.Println(plaintext)

	// 加密
	ciphertext, err := utils.AesEncrypt(plaintext, key)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return
	}

	fmt.Printf("Encrypted: %x\n", ciphertext)

	// 解密
	decryptedText, err := utils.AesDecrypt(ciphertext, key)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}

	fmt.Println("Decrypted:", string(decryptedText))
*/
