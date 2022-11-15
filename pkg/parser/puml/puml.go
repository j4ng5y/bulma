package puml

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog/log"
)

type (
	Parser struct {
		f string
	}
)

const (
	pumlStart = "@startuml"
	pumlEnd   = "@enduml"
)

func NewParser(file string) (*Parser, error) {
	s, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		return nil, fmt.Errorf("%s is a directory, not a file", file)
	}
	return &Parser{f: file}, nil
}

func (p *Parser) countLines(r io.Reader) (int, error) {
	var (
		count int
		read  int
		pos   int
		err   error
		nl    byte = '\n'
	)

	// A basic buffer to read 32Kb chunks of the file into
	buf := make([]byte, 32*1024)

	for {
		// Read the file contents into our buffer up to 32Kb
		read, err = r.Read(buf)
		if err != nil {
			break
		}

		for {
			// Check if our newline character exists in the read in bytes
			idxOf := bytes.IndexByte(buf[pos:read], nl)
			if idxOf == -1 {
				// If we have reached the end of the file, and there are still contents in the buffer, assume we have
				// a final line without a newline character
				if len(buf[pos:read]) > 0 {
					count++
				}
				// Otherwise, break out of the loop
				break
			}

			// Add to the running count of lines
			count++
			// Shift the buffer to look for more newlines
			pos += idxOf + 1
		}
	}

	if err == io.EOF {
		// If we get to the end of the file, return the count
		return count, nil
	}
	// Otherwise, return the error we encountered
	return count, err
}

func (p *Parser) Parse() error {
	log.Debug().Msgf("parsing %s", p.f)
	f, err := os.Open(p.f)
	if err != nil {
		return fmt.Errorf("PUML.Parser.Parse: unable to open file: %s, err: %w", p.f, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Warn().Err(err).Send()
		}
	}()

	l, err := p.countLines(f)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	line := 0
	for scanner.Scan() {
		line++
		if line == 1 {
			if scanner.Text() != pumlStart {
				return fmt.Errorf("file does not begin with %s", pumlStart)
			}
		}
		if line == l {
			if scanner.Text() != pumlEnd {
				return fmt.Errorf("file does not end with %s", pumlEnd)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
