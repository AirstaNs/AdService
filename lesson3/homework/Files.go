package main

import "bufio"

type TransformReader struct {
	reader bufio.Reader
	offset int
	limit  int64
}

type TransformWriter struct {
	writer      bufio.Writer
	Conversions Converter
}

type Input interface {
	Read(buff []byte) (n int, err error)
	Discard() (discarded int, err error)
}
type Writer interface {
	Write(buff []byte) (n int, err error)
}

func (trRead *TransformReader) Read(buff []byte) (n int, err error) {
	return trRead.reader.Read(buff)
}

func (trRead *TransformReader) Discard() (discarded int, err error) {
	return trRead.reader.Discard(trRead.offset) // TODO
}

func (trWrite *TransformWriter) Write(buff []byte) (n int, err error) {
	nn, err := trWrite.writer.Write(buff)
	trWrite.writer.Flush()
	return nn, err
}
