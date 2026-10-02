package utils

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/duke-git/lancet/v2/validator"
	"github.com/maiguangyang/aescrypto"
)

// aes密钥
var secretKey string = "dAcGmFaAgBaAnHbOcLdIbCiUbIuAfFbA"

// aes解密
func AesLeftDecrypt(text string, key string) (string, error) {
	if validator.IsEmptyString(key) {
		key = secretKey
	}
	// 先base64解密
	prvPem, err := base64.StdEncoding.DecodeString(text)

	// 再使用hex.DecodeString转换
	dexDecode, err := hex.DecodeString(string(prvPem))

	if err != nil {
		return "", err
	}

	// 最后解密对称加密的内容
	prvPem, err = AesDecrypt(dexDecode, []byte(key))
	if err != nil {
		return "", err
	}

	return string(prvPem), nil
}

// AES CBC加密
func AesCbcEncrypt(text string, key string) string {
	if validator.IsEmptyString(key) {
		key = secretKey
	}
	cryptText, err := aescrypto.AesCbcPkcs7Encrypt([]byte(text), []byte(key), nil)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(cryptText)
}

// AES CBC解密
func AesCbcDecrypt(text string, key string) string {
	if validator.IsEmptyString(key) {
		key = secretKey
	}

	cryptText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return ""
	}

	decData, err := aescrypto.AesCbcPkcs7Decrypt(cryptText, []byte(key), nil)

	if err != nil {
		return ""
	}

	return string(decData)
}
