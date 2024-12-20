package mysql_protocol

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"math/rand"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func parse3BytesToUint32(data []byte) uint32 {
	if len(data) != 3 {
		panic("Expected 3 bytes for parsing")
	}
	return uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16
}

func scramblePassword(password string, scramble []byte) []byte {
	hash1 := sha1.Sum([]byte(password))
	hash2 := sha1.Sum(hash1[:])
	scrambleHash := sha1.New()
	scrambleHash.Write(scramble)
	scrambleHash.Write(hash2[:])
	finalHash := scrambleHash.Sum(nil)

	result := make([]byte, len(finalHash))
	for i := range finalHash {
		result[i] = finalHash[i] ^ hash1[i]
	}
	return result
}

func readNullTerminatedString(reader *bytes.Reader) (string, error) {
	var strBytes []byte
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return "", err
		}
		if b == 0 {
			break
		}
		strBytes = append(strBytes, b)
	}
	return string(strBytes), nil
}

func readFixedLengthString(rw *bytes.Reader, length int) (string, error) {
	var data = make([]byte, length)
	var _, err = rw.Read(data)

	if err != nil {
		return "", err
	}

	if len(data) < length {
		return "", errors.New("not enough bytes to read fixed-length string")
	}

	// Extract the fixed-length string and the remaining bytes
	str := string(data)

	return str, nil
}

func PrintBytesBinary(biter []byte) {
	for _, b := range biter {
		fmt.Printf("0x%x,", b)
	}

	fmt.Printf("\n")
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ&*^")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
