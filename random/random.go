package random

import (
	"bytes"
	cryptorand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"io"
	"math/rand/v2"
	"strings"
)

var (
	numeric      string = "0123456789"                         // digits: [0-9]
	lowercase    string = "abcdefghijklmnopqrstuvwxyz"         // asci lowerrcase letters: [a-z]
	uppercase    string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"         // ascii ippercase letters: [A-Z]
	alpha        string = lowercase + uppercase                // ascii letters: [a-zA-Z]
	alphanumeric string = alpha + numeric                      // ascii charaters: [a-zA-Z0-9]
	special      string = "~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`" // special characters
	ascii        string = alphanumeric + special               // all ascii characters
)

func rndByte(n int) []byte {
	b := make([]byte, n)
	_, err := cryptorand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

// Alpha generates a random alphabet string of specified length from letters [a-zA-Z]
func Alpha(length int) string {
	runes := []rune(alpha)

	var bb bytes.Buffer
	bb.Grow(length)
	l := uint32(len(runes))

	for range length {
		bb.WriteRune(runes[binary.BigEndian.Uint32(rndByte(4))%l])
	}
	return bb.String()
}

// AlphaNumeric generates a random alphanumeric string of specified length from letters [a-zA-Z0-9]
func AlphaNumeric(length int) string {
	runes := []rune(alpha)

	var bb bytes.Buffer
	bb.Grow(length)
	l := uint32(len(runes))

	for range length {
		bb.WriteRune(runes[binary.BigEndian.Uint32(rndByte(4))%l])
	}
	return bb.String()
}

// AlphaLowercase generates a random alphabet lowercase string of specified length from letters [a-z]
func AlphaLowercase(length int) string {
	runes := []rune(lowercase)

	var bb bytes.Buffer
	bb.Grow(length)
	l := uint32(len(runes))

	for range length {
		bb.WriteRune(runes[binary.BigEndian.Uint32(rndByte(4))%l])
	}
	return bb.String()
}

// AlphaUppercase generates a random alphabet uppercase string of specified length from letters [A-Z]
func AlphaUppercase(length int) string {
	runes := []rune(uppercase)

	var bb bytes.Buffer
	bb.Grow(length)
	l := uint32(len(runes))

	for range length {
		bb.WriteRune(runes[binary.BigEndian.Uint32(rndByte(4))%l])
	}
	return bb.String()
}

// Ascii generates a random ascii string of specified length
func Ascii(length int) string {
	runes := []rune(ascii)

	var bb bytes.Buffer
	bb.Grow(length)
	l := uint32(len(runes))

	for range length {
		bb.WriteRune(runes[binary.BigEndian.Uint32(rndByte(4))%l])
	}
	return bb.String()
}

// Token generates a random token
func Token(length int) string {
	b := make([]byte, length)
	if _, err := io.ReadFull(cryptorand.Reader, b); err != nil {
		panic(err.Error())
	}

	token := base64.URLEncoding.EncodeToString(b)

	withoutPadding := strings.TrimRight(token, "=")
	return withoutPadding
}

// AlphaInRange generates a random alphabet string in range
func AlphaInRange(min, max int) string {
	n := rand.IntN((max-min)+1) + min
	return Alpha(n)
}

// AlphaNumericInRange generates a random alphanumeric string in range
func AlphaNumericInRange(min, max int) string {
	n := rand.IntN((max-min)+1) + min
	return AlphaNumeric(n)
}

// AlphaLowercaseInRange generates a random alphabet lowercase string in range
func AlphaLowercaseInRange(min, max int) string {
	n := rand.IntN((max-min)+1) + min
	return AlphaLowercase(n)
}

// AlphaUppercaseInRange generates a random alphabet uppercase string in range
func AlphaUppercaseInRange(min, max int) string {
	n := rand.IntN((max-min)+1) + min
	return AlphaUppercase(n)
}

// AsciiInRange generates a random ascii string in range
func AsciiInRange(min, max int) string {
	n := rand.IntN((max-min)+1) + min
	return Ascii(n)
}

// IntInRange generates a random int in range
func IntInRange(min, max int) int {
	return rand.IntN((max-min)+1) + min
}

// Bool generates a random bool
func Bool() bool {
	return rand.IntN(2) == 1
}
