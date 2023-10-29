package main

import "fmt"

type ReaderSim struct {
	Content      []byte
	CurrentIndex int
}

func (r *ReaderSim) Read(buf []byte) (int, error) {
	capacity := cap(buf)
	readStart := r.CurrentIndex
	if readStart >= len(r.Content) {
		return 0, fmt.Errorf("reader index is passed the end")
	}
	readEnd := readStart + capacity
	if readEnd >= len(r.Content) {
		readEnd = len(r.Content)
	}
	readLen := readEnd - readStart
	copy(buf, r.Content[readStart:readEnd])
	r.CurrentIndex += readLen
	return readLen, nil
}
