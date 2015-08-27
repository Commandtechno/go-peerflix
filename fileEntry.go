package main

import (
	"io"
	"os"

	"github.com/anacrolix/torrent"
)

// SeekableContent describes an io.ReadSeeker that can be closed as well.
type SeekableContent interface {
	io.ReadSeeker
	io.Closer
}

// FileEntry helps reading a torrent file.
type FileEntry struct {
	File   *torrent.File
	Reader *torrent.Reader
}

func (f FileEntry) Read(p []byte) (n int, err error) {
	return f.Reader.Read(p)
}

func (f FileEntry) Seek(offset int64, whence int) (int64, error) {
	return f.Reader.Seek(offset+f.File.Offset(), whence)
}

func (f FileEntry) Close() error {
	return f.Reader.Close()
}

// NewFileReader sets up a torrent file for streaming reading.
func NewFileReader(f torrent.File) SeekableContent {
	// We read ahead 1% of the file continuously.
	var readahead = f.Length() / 100

	// We begin by prioritizing 5% of the beginning of the file.
	f.PrioritizeRegion(f.Offset(), readahead*5)

	reader := t.NewReader()
	reader.SetReadahead(readahead)
	reader.SetResponsive()
	reader.Seek(f.Offset(), os.SEEK_SET)

	return &FileEntry{
		File:   &f,
		Reader: reader,
	}
}