package main

import (
	"flag"
	"log"
	"math"
	"strconv"
)

type Options struct {
	From        string
	To          string
	Offset      int64
	Limit       int64
	BlockSize   int64
	Conversions Converter
}

const (
	input     = "stdin"
	output    = "stdout"
	offset    = 0
	limit     = math.MaxInt32
	blockSize = 4096
)

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", input, "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", output, "file to write. by default - stdout")
	off := flag.String("offset", strconv.Itoa(offset), "number of bytes inside input to be skipped when copying")
	lim := flag.String("limit", strconv.Itoa(limit), "max number of bytes to read")
	bs := flag.String("block-size", strconv.Itoa(blockSize), "size of one block in bytes")
	conv := flag.String("conv", DefaultConv, "comma-separated list of conversion functions")
	flag.Parse()
	opts.Offset = parseIntFlag(*off)
	opts.Limit = parseIntFlag(*lim)
	opts.BlockSize = parseIntFlag(*bs)

	var opt = new(TransformOptions)
	opt.parse(conv)
	opts.Conversions = opt
	return &opts, nil
}
func parseIntFlag(value string) int64 {
	arg, err := strconv.ParseInt(value, 10, 64)
	handleAllError(err, "can not parse int flag")

	if arg < 0 {
		log.Fatal("negative value")
	}
	return arg
}
