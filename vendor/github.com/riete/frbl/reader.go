package frbl

import (
	"bufio"
	"io"
	"os"
)

type FileReader interface {
	ReadLine() error
	Content() chan string
	Close()
}

type file struct {
	path    string
	offset  int64
	file    *os.File
	content chan string
}

func NewFileReader(path string) FileReader {
	content := make(chan string)
	return &file{path: path, offset: offsetGet(path), content: content}
}

func (f *file) open() error {
	var err error
	if f.file, err = os.Open(f.path); err != nil {
		return err
	}
	return nil
}

func (f *file) forRotated() error {
	if f.offset == 0 {
		return nil
	}
	if end, err := f.file.Seek(0, io.SeekEnd); err != nil {
		return err
	} else {
		if f.offset > end {
			f.offset = 0
		}
		return nil
	}
}

func (f *file) setOffset() error {
	if offset, err := f.file.Seek(0, io.SeekCurrent); err != nil {
		return err
	} else {
		f.offset = offset
		return offsetUpdate(f.path, offset)
	}
}

func (f *file) seek() error {
	_, err := f.file.Seek(f.offset, io.SeekStart)
	return err

}

func (f *file) ReadLine() error {
	if err := f.open(); err != nil {
		return err
	}
	defer f.file.Close()
	if err := f.forRotated(); err != nil {
		return err
	}
	if err := f.seek(); err != nil {
		return err
	}
	r := bufio.NewReader(f.file)
	for {
		data, err := r.ReadBytes('\n')
		if err == io.EOF {
			return f.setOffset()
		}
		f.content <- string(data)
	}
}

func (f file) Content() chan string {
	return f.content
}

func (f *file) Close() {
	close(f.content)
}
