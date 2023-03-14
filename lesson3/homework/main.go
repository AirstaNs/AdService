package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stderr) //TODO
	opts, err := ParseFlags()
	handleAllError(err, "can not parse flags:")

	var read, write *os.File
	read = getInput(opts)
	defer read.Close()
	handleAllError(err, "can not open file:")

	write = getOutput(opts)
	defer write.Close()
	handleAllError(err, "can not open file:")

	// seek to offset
	//	n, err := read.Seek((opts.Offset), io.SeekStart)
	_, _ = io.CopyN(io.Discard, read, opts.Offset)
	_, err = io.CopyN(write, read, opts.Limit)

	if err != nil {
		if errors.Is(err, io.EOF) {
			fmt.Print("")
			//fmt.Print(io.EOF)
			//fmt.Println(io.EOF)
		} else {
			log.Fatal(err)
		}
	}
}

func openFile(path string, opts *Options) *os.File {
	log.SetOutput(os.Stderr) //TODO
	file, err := os.Open(path)
	handleAllError(err, "can not open file:")
	return file
}
func createFile(path string) *os.File {
	log.SetOutput(os.Stderr) //TODO
	file, err := os.Open(path)
	if err != nil {
		file, err = os.Create(path)
		handleAllError(err, "can not create file:")
	} else {
		log.Fatal("file already exists")
	}
	return file

}

func handleAllError(err error, message string) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func writeFile(opts *Options, inputFile *os.File, outputFile *os.File) {
	buf := make([]byte, opts.BlockSize)
	for {
		readTotal, err := inputFile.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("EOF")
				break // after reading the last chunk, break the loop
			}
			log.Fatal(err)
		}
		outputFile.WriteString(string(buf[:readTotal]))
		//fmt.Println(string(buf[:readTotal]))
	}
}
func readFile(opts *Options, inputFile *os.File) {
	buf := make([]byte, opts.BlockSize)
	for {
		readTotal, err := inputFile.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("EOF")
				break // after reading the last chunk, break the loop
			}
			log.Fatal(err)
		}
		fmt.Println(string(buf[:readTotal]))
	}
}

func getInput(opts *Options) *os.File {
	if opts.From == input {
		return os.Stdin
	}
	return openFile(opts.From, opts)
}
func getOutput(opts *Options) *os.File {
	if opts.To == output {
		return os.Stdout
	}
	return createFile(opts.To)
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
