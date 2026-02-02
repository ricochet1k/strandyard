package activity

import (
	"bytes"
	"io"
	"os"
)

// ReverseScanner scans a file backwards line by line.
type ReverseScanner struct {
	file   *os.File
	pos    int64
	buffer []byte
	err    error
	line   []byte
}

// NewReverseScanner creates a new ReverseScanner for the given file.
func NewReverseScanner(file *os.File) (*ReverseScanner, error) {
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return &ReverseScanner{
		file: file,
		pos:  info.Size(),
	}, nil
}

// Scan moves to the previous line. Returns false if start of file is reached or an error occurs.
func (s *ReverseScanner) Scan() bool {
	if s.err != nil {
		return false
	}

	const blockSize = 4096
	for {
		// Look for a newline in the current buffer (from the end)
		if len(s.buffer) > 0 {
			idx := bytes.LastIndexByte(s.buffer, '\n')
			if idx >= 0 {
				s.line = s.buffer[idx+1:]
				s.buffer = s.buffer[:idx]
				return true
			}
		}

		// If no newline and we are at the start of the file, the rest of the buffer is the first line
		if s.pos <= 0 {
			if len(s.buffer) > 0 {
				s.line = s.buffer
				s.buffer = nil
				return true
			}
			return false
		}

		// Read the next block from the end
		readSize := int64(blockSize)
		if s.pos < readSize {
			readSize = s.pos
		}
		s.pos -= readSize

		newBuf := make([]byte, readSize)
		_, err := s.file.ReadAt(newBuf, s.pos)
		if err != nil && err != io.EOF {
			s.err = err
			return false
		}

		// Prepend the new block to the buffer
		s.buffer = append(newBuf, s.buffer...)
	}
}

// Text returns the current line as a string.
func (s *ReverseScanner) Text() string {
	return string(s.line)
}

// Err returns the first non-EOF error encountered.
func (s *ReverseScanner) Err() error {
	return s.err
}
