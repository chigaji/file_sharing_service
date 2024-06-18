package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

var key = []byte("e40adb88-d03b-48ce-a2ea-374616eb") // 32 bytes

func Encrypt(input io.Reader, output io.Writer) error {

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, aes.BlockSize)

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	if _, err := output.Write(iv); err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	writer := &cipher.StreamWriter{S: stream, W: output}

	if _, err := io.Copy(writer, input); err != nil {
		return err
	}
	return nil
}

func Decrypt(input io.Reader) io.Reader {

	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err)
	}

	iv := make([]byte, aes.BlockSize)

	if _, err := io.ReadFull(input, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	return &cipher.StreamReader{S: stream, R: input}
}
