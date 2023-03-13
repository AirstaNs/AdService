package main

import (
	"flag"
	"fmt"
	"math"
	"os"
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

func main() {
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}
	fmt.Println(opts.Offset, opts.To)
}

// -conv upper_case,lower_case,trim_spaces
//func main() {
//	opts, err := ParseFlags()
//	if err != nil {
//		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
//		os.Exit(1)
//	}
//	var reader *bufio.Reader
//	var writer *bufio.Writer
//	if opts.To == "stdout" {
//		writer = bufio.NewWriter(os.Stdout)
//	} else {
//		wr, err := os.OpenFile(opts.From, os.O_APPEND|os.O_CREATE, 0644)
//		if err != nil {
//			os.Stderr.WriteString(err.Error())
//			os.Exit(1)
//		}
//		writer = bufio.NewWriter(wr)
//		defer wr.Close()
//	}
//	if opts.From == "stdin" {
//		reader = bufio.NewReader(os.Stdin)
//	} else {
//		read, err := os.OpenFile(opts.From, os.O_RDONLY|os.O_CREATE, 0644)
//		if err != nil {
//			fmt.Println(err)
//			os.Exit(1)
//		}
//		reader = bufio.NewReader(read)
//		defer read.Close()
//	}
//
//	var input Input = &TransformReader{reader: *reader, offset: opts.Offset, limit: opts.Limit}
//	var output Writer = &TransformWriter{writer: *writer, Conversions: opts.Conversions}
//	//writer.WriteString("aaaa")
//	//writer.Flush()
//	input.Discard() // TODO
//	//reader.Discard(opts.Offset) // -offset
//	data := make([]byte, opts.BlockSize)
//	tempData := make([]byte, opts.BlockSize) // -block-size
//	for {
//		//	fmt.Println(reader.Buffered(), " buff")
//		n, err := input.Read(tempData)
//		//fmt.Println(reader.Buffered())
//		//writer.Write(tempData[:n])
//		data = append(data, tempData[:n]...)
//		if err == io.EOF { // если конец файла
//			break // выходим из цикла
//		}
//	}
//	convert := opts.Conversions.Convert(data)
//	output.Write(convert)
//	//writer.Write(convert)
//	//writer.Flush()
//}
