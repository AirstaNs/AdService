package main

import (
	"flag"
	"math"
	"strconv"
)

type Options struct {
	From        string
	To          string
	Offset      int
	Limit       int
	BlockSize   int
	Conversions Converter
}

const (
	input     = "stdin"
	output    = "stdout"
	offset    = 0
	limit     = math.MaxInt32
	blockSize = 4096
	conv      = ""
)

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", input, "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", output, "file to write. by default - stdout")
	off := flag.String("offset", strconv.Itoa(offset), "number of bytes inside input to be skipped when copying")
	lim := flag.String("limit", strconv.Itoa(limit), "max number of bytes to read")
	bs := flag.String("block-size", strconv.Itoa(blockSize), "size of one block in bytes")
	conv := flag.String("conv", conv, "comma-separated list of conversion functions")
	flag.Parse()
	opts.Offset = parseIntFlag(*off, offset)
	opts.Limit = parseIntFlag(*lim, limit)
	opts.BlockSize = parseIntFlag(*bs, blockSize)

	var opt = new(TransformOptions)
	opt.parse(conv)
	opts.Conversions = opt
	return &opts, nil
}
func parseIntFlag(value string, defaultArg int) int {
	arg, err := strconv.Atoi(value)
	if err != nil {
		arg = defaultArg
	}
	return arg
}
