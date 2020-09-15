package mobTech

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"log"
	"fmt"
	"sort"
	"errors"
)

type base struct {}

func (m *base) PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func (m *base) PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 数据加密
func (m *base) DesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	origData = m.PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, []byte("00000000"))
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// 数据解密
func (m *base) DesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, []byte("00000000"))
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = m.PKCS5UnPadding(origData)
	return origData, nil
}

// base64编码
func (m *base) Base64Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}

// base64解码
func (m *base) Base64Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}

// 生成签名
func (m *base) generateSign(request map[string]interface{}, secret string) string {
	ret := ""
	var keys []string
	for k := range request {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, k := range keys {
		ret = ret + fmt.Sprintf("%v=%v&", k, request[k])
	}
	ret = ret[:len(ret)-1] + secret

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(ret))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// post请求
func (m *base) HttpPostBody(url string, msg map[string]interface{}) ([]byte, error) {
	request := gorequest.New()
	resp, body, errs := request.Post(url).Set("Content-Type", "application/json; charset=utf-8").SendMap(msg).End()
	if errs != nil {
		log.Fatalf("mob_trace: request err:%+v, body:%s", errs, body)
		return nil, errors.New("request error")
	}

	return ioutil.ReadAll(resp.Body)
}

