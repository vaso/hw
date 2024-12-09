package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if input == "" {
		return "", nil
	}
	var stringChunks []string
	stringToPrint := ""
	isEscaped := false
	for _, sym := range input {
		isDigit := unicode.IsDigit(sym)
		strSym := string(sym)
		switch {
		case isEscaped:
			if !isDigit && strSym != "\\" {
				return "", ErrInvalidString
			}
			stringToPrint = strSym
			isEscaped = false
		case isDigit:
			res, err := repeatString(stringToPrint, sym)
			if err != nil {
				return "", ErrInvalidString
			}
			stringChunks = append(stringChunks, res)
			stringToPrint = ""
		case strSym == "\\":
			isEscaped = true
			stringChunks = append(stringChunks, stringToPrint)
		case stringToPrint == "":
			stringToPrint = strSym
		default:
			stringChunks = append(stringChunks, stringToPrint)
			stringToPrint = strSym
		}
	}
	if stringToPrint != "" {
		stringChunks = append(stringChunks, stringToPrint)
	}

	return strings.Join(stringChunks, ""), nil
}

func repeatString(str string, counterSymbol rune) (string, error) {
	if str == "" {
		return "", ErrInvalidString
	}
	if !unicode.IsDigit(counterSymbol) {
		return "", ErrInvalidString
	}
	counter, err := strconv.Atoi(string(counterSymbol))
	if err != nil {
		return "", err
	}
	if counter == 0 {
		return "", nil
	}

	return strings.Repeat(str, counter), nil
}
