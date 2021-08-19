package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	//"crypto/sha256"
	//"io"
	"encoding/base64"
	"strings"
)

func packageData(originalData []byte, packageSize int) (r [][]byte) {
	var src = make([]byte, len(originalData))
	copy(src, originalData)

	r = make([][]byte, 0)
	if len(src) <= packageSize {
		return append(r, src)
	}
	for len(src) > 0 {
		var p = src[:packageSize]
		r = append(r, p)
		src = src[packageSize:]
		if len(src) <= packageSize {
			r = append(r, src)
			break
		}
	}
	return r
}

func RSAEncrypt(plaintext, key []byte) ([]byte, error) {
	var err error
	var block *pem.Block
	block, _ = pem.Decode(key)
	if block == nil {
		return nil, errors.New("public key error")
	}

	var pubInterface interface{}
	pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	var pub = pubInterface.(*rsa.PublicKey)

	var data = packageData(plaintext, pub.N.BitLen()/8-11)
	var cipherData []byte = make([]byte, 0, 0)

	for _, d := range data {
		var c, e = rsa.EncryptPKCS1v15(rand.Reader, pub, d)
		if e != nil {
			return nil, e
		}
		cipherData = append(cipherData, c...)
	}

	return cipherData, nil
}

func RSADecrypt(ciphertext, key []byte) ([]byte, error) {
	var err error
	var block *pem.Block
	block, _ = pem.Decode(key)
	if block == nil {
		return nil, errors.New("private key error")
	}

	var pri *rsa.PrivateKey
	pri, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	var data = packageData(ciphertext, pri.PublicKey.N.BitLen()/8)
	var plainData []byte = make([]byte, 0, 0)

	for _, d := range data {
		var p, e = rsa.DecryptPKCS1v15(rand.Reader, pri, d)
		if e != nil {
			return nil, e
		}
		plainData = append(plainData, p...)
	}
	return plainData, nil
}

func SignPKCS1v15(src, key []byte, hash crypto.Hash) ([]byte, error) {
	var h = hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)

	var err error
	var block *pem.Block
	block, _ = pem.Decode(key)
	if block == nil {
		return nil, errors.New("private key error")
	}

	var pri *rsa.PrivateKey
	pri, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.SignPKCS1v15(rand.Reader, pri, hash, hashed)
}

func VerifyPKCS1v15(src, sign, key string, hash crypto.Hash) error {
	sign = strings.Replace(sign, " ", "+", -1)
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return  err
	}

	var h = hash.New()
	h.Write([]byte(src))
	var hashed = h.Sum(nil)

	//var block *pem.Block
	//block, _ = pem.Decode(key)
	//if block == nil {
	//	model.Log.Error("public key err:%v, key:%v", err, string(key))
	//}

	buf, err := base64.StdEncoding.DecodeString(key)
	if buf == nil {
		return errors.New("failed to parse PEM block containing the public key")
	}

	// 计算代签名字串的SHA256哈希
	//t := sha256.New()
	//io.WriteString(t, string(src))
	//digest := t.Sum(nil)

	var pubInterface interface{}
	pubInterface, err = x509.ParsePKIXPublicKey(buf)
	if err != nil {
		return err
	}
	var pub = pubInterface.(*rsa.PublicKey)

	return rsa.VerifyPKCS1v15(pub, hash, hashed, signBytes)
}
