package record

import (
	"bufio"
	"context"
	"io"
	"net"
	"os"
	"time"
)

type Playback struct {
	connection net.Conn
	file       *os.File
}

// TODO rewite to make it use io.Reader Writer
func OpenPlayback(address string, filename string) (*Playback, error) {
	c, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return &Playback{
		connection: c,
		file:       f,
	}, nil
}

// Run reads data from the file and writes it to the host.
func (p *Playback) Run(ctx context.Context) error {
	pt := 0

	r := bufio.NewReader(p.file)
	for ctx.Err() == nil {
		line, err := r.ReadBytes(byte('\n'))
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// remove trailing whitespace
		j := len(line) - 1
		for j >= 0 && line[j] <= ' ' {
			j--
		}
		line = line[:j+1]

		// skip line if empty or comment
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		dt := 200 * time.Millisecond
		if '0' <= line[0] && line[0] <= '9' {
			// line starts with timestamp in mS
			t := 0
			i := 0
			for line[i] != ' ' {
				t *= 10
				t += int(line[i] - '0')
				i++
			}
			for line[i] == ' ' {
				i++
			}
			line = line[i:]

			dt = time.Duration(t-pt) * time.Millisecond
			pt = t

			if dt > 5*time.Second {
				dt = 5 * time.Second
			}
		}
		time.Sleep(dt)

		_, err = p.connection.Write(line)
		if err != nil {
			return err
		}
		_, err = p.connection.Write([]byte{'\r', '\n'})
		if err != nil {
			return err
		}
	}

	return nil
}

// Close closes the playback.
func (rr *Playback) Close() error {
	err1 := rr.file.Close()
	err2 := rr.connection.Close()
	if err1 != nil {
		return err1
	}
	return err2
}
