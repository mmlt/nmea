package record

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"time"
)

type Record struct {
	connection net.Conn
	file       *os.File
	timestamp  bool
}

//
func Open(address string, filename string, timestamp bool) (*Record, error) {
	c, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return &Record{
		connection: c,
		file:       f,
		timestamp:  timestamp,
	}, nil
}

// Run reads data from the host and writes it to the file.
func (rr *Record) Run(ctx context.Context) error {
	r := bufio.NewReader(rr.connection)
	for ctx.Err() == nil {
		line, err := r.ReadBytes(byte('\n'))
		if err != nil {
			return err
		}

		// remove trailing whitespace
		j := len(line) - 1
		for line[j] <= ' ' && j > 0 {
			j--
		}
		line = line[:j+1]

		if rr.timestamp {
			t := time.Now().UnixMilli()
			fmt.Fprintf(rr.file, "%d ", t)
		}

		_, err = rr.file.Write(line)
		if err != nil {
			return err
		}
		_, err = rr.file.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// Close closes the recording.
func (rr *Record) Close() error {
	err1 := rr.file.Close()
	err2 := rr.connection.Close()
	if err1 != nil {
		return err1
	}
	return err2
}
