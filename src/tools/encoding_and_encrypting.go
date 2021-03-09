package tools

import (
	"encoding/base64"
)

func Encrypt(text string) string {
	nonce := make([]byte, Config.Gcm.NonceSize())
	return B64Enc(string(Config.Gcm.Seal(nonce, nonce, []byte(text), nil)))
}

func Decrypt(srcPas string) (string, error) {
	srcPas, err := B64Dec(srcPas)
	if err != nil {
		return "", err
	}
	pas := []byte(srcPas)
	nonceSize := Config.Gcm.NonceSize()
	nonce, pas := pas[:nonceSize], pas[nonceSize:]
	text, err := Config.Gcm.Open(nil, nonce, pas, nil)
	if err != nil {
		return "", err
	}
	return string(text), nil

}

func B64Enc(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))

}

func B64Dec(text string) (string, error) {
	result, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
