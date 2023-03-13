package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
)

type Options struct {
	From        string
	To          string
	Offset      int
	Limit       int64
	BlockSize   int64
	Conversions Converter
}

func ParseFlags() (*Options, error) {
	var opts Options
	flag.StringVar(&opts.From, "from", "stdin", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "stdout", "file to write. by default - stdout")
	flag.IntVar(&opts.Offset, "offset", 0, "number of bytes inside input to be skipped when copying")
	flag.Int64Var(&opts.Limit, "limit", math.MaxInt32, "max number of bytes to read")
	flag.Int64Var(&opts.BlockSize, "block-size", 4096, "size of one block in bytes")
	//flag.Args(&opts.Conversions, "conv", []string{}, "size of one block in bytes")
	//flag.Var(&opts.Conversions, "conv", "comma-separated list of conversion functions")
	// todo: parse and validate all flags
	conv := flag.String("conv", "", "comma-separated list of conversion functions")
	flag.Parse()

	var opt = new(TransformOptions)
	opt.parse(conv)
	opts.Conversions = opt
	return &opts, nil
}

// -conv upper_case,lower_case,trim_spaces
func main() {
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}
	var reader *bufio.Reader
	var writer *bufio.Writer
	if opts.To == "stdout" {
		writer = bufio.NewWriter(os.Stdout)
	} else {
		wr, err := os.OpenFile(opts.From, os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
		writer = bufio.NewWriter(wr)
		defer wr.Close()
	}
	if opts.From == "stdin" {
		reader = bufio.NewReader(os.Stdin)
	} else {
		read, err := os.OpenFile(opts.From, os.O_RDONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		reader = bufio.NewReader(read)
		defer read.Close()
	}

	var input Input = &TransformReader{reader: *reader, offset: opts.Offset, limit: opts.Limit}
	var output Writer = &TransformWriter{writer: *writer, Conversions: opts.Conversions}
	//writer.WriteString("aaaa")
	//writer.Flush()
	input.Discard() // TODO
	//reader.Discard(opts.Offset) // -offset
	data := make([]byte, opts.BlockSize)
	tempData := make([]byte, opts.BlockSize) // -block-size
	for {
		//	fmt.Println(reader.Buffered(), " buff")
		n, err := input.Read(tempData)
		//fmt.Println(reader.Buffered())
		//writer.Write(tempData[:n])
		data = append(data, tempData[:n]...)
		if err == io.EOF { // если конец файла
			break // выходим из цикла
		}
	}
	convert := opts.Conversions.Convert(data)
	output.Write(convert)
	//writer.Write(convert)
	//writer.Flush()
}
