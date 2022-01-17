package thirdLogin

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func AesDecrypt(encryptedData, sessionKey, iv string) ([]byte, error) {
	// Base64解码
	keyBytes, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, err
	}
	
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}
	
	cryptData, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	
	origData := make([]byte, len(cryptData))
	// AES
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}
	
	// CBC
	mode := cipher.NewCBCDecrypter(block, ivBytes)
	// 解密
	mode.CryptBlocks(origData, cryptData)
	// 去除填充位
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	if length > 0 {
		unPadding := int(plantText[length-1])
		return plantText[:(length - unPadding)]
	}
	
	return plantText
}
