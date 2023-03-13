package main

import (
	"bytes"
	"strings"
)

const (
	Upper      = "upper_case"
	Lower      = "lower_case"
	TrimSpaces = "trim_spaces"
)

type TransformOptions struct {
	Upper, Lower, Trim bool
}

type Converter interface {
	Convert(bytes []byte) []byte
}

func (opt *TransformOptions) Convert(arrByte []byte) []byte {
	if opt.Trim {
		arrByte = bytes.TrimSpace(arrByte)
	}
	if opt.Lower {
		arrByte = bytes.ToLower(arrByte)
	}
	if opt.Upper {
		arrByte = bytes.ToUpper(arrByte)
	}
	return arrByte
}
func (opt *TransformOptions) parse(conv *string) {
	args := parsing(*conv)

	for _, arg := range args {
		switch arg {
		case Upper:
			opt.Upper = true
		case Lower:
			opt.Lower = true
		case TrimSpaces:
			opt.Trim = true
		}
	}
}

func parsing(str string) []string {
	const delimiter = ","
	return strings.Split(str, delimiter)
}
