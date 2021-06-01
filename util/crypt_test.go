package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	expected := "asd"

	resultEncrypt := Encrypt("asd")
	resultDecrypt := Decrypt(resultEncrypt)

	assert.NotEqual(t, expected, resultEncrypt, "Deve encriptar o texto")
	assert.Equal(t, expected, resultDecrypt, "Deve desencriptar o texto")
}

func TestBase64EncodeDecode(t *testing.T) {
	expected := "asd"

	resultEncode := Base64("asd")
	resultDecode := Base64Decode(resultEncode)

	assert.NotEqual(t, expected, resultEncode, "Deve encodar em Base64 o texto")
	assert.Equal(t, expected, resultDecode, "Deve desencodar o texto Base64")
}
