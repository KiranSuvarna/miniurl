package misc

import (
	"errors"
	"math"
	"strings"
)

var (
	Base         = 62
	CharacterSet = "K7LawG8eGY46BjvdJWhaznN2YoPNd0Va3xgLKgEdOmkugYtpRwcq3M3dP6GcAErNoQu4xu"
)

func Encode(num int) string {
	b := make([]byte, 0)

	// loop as long the num is bigger than zero
	for num > 0 {
		// receive the rest
		r := math.Mod(float64(num), float64(Base))

		// devide by Base
		num /= Base

		// append chars
		b = append([]byte{CharacterSet[int(r)]}, b...)
	}

	return string(b)
}

func Decode(s string) (int, error) {
	var r, pow int

	// loop through the input
	for i, v := range s {
		// convert position to power
		pow = len(s) - (i + 1)

		// IndexRune returns -1 if v is not part of CharacterSet.
		pos := strings.IndexRune(CharacterSet, v)

		if pos == -1 {
			return pos, errors.New("invalid character: " + string(v))
		}

		// calculate
		r += pos * int(math.Pow(float64(Base), float64(pow)))
	}

	return int(r), nil
}
