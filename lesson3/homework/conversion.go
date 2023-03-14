package main

import (
	"bytes"
	"log"
	"strings"
)

const (
	Upper       = "upper_case"
	Lower       = "lower_case"
	TrimSpaces  = "trim_spaces"
	DefaultConv = "defaultConv"
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
	args := splitArgs(*conv)

	for _, arg := range args {
		switch arg {
		case DefaultConv:
			opt.Upper, opt.Lower, opt.Lower = false, false, false
		case Upper:
			opt.Upper = true
		case Lower:
			opt.Lower = true
		case TrimSpaces:
			opt.Trim = true
		default:
			log.Fatal("unknown conversion function")
		}
	}

	checkContradictory(opt)

}
func checkContradictory(opt *TransformOptions) {
	if opt.Upper && opt.Lower {
		log.Fatal("upper and lower case can not be used together")
	}
}

func splitArgs(str string) []string {
	const delimiter = ","
	return strings.Split(str, delimiter)
}
