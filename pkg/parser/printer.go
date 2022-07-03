package parser

import (
	"bytes"
	"fmt"
)

func Print(s Sentence) (string, error) {
	w := &bytes.Buffer{}
	fmt.Fprint(w, "$", s.TalkerID(), s.DataType())

	p := printers[s.DataType()]
	if p == nil {
		return "", fmt.Errorf("no printer for: %s", s.DataType())
	}

	err := p(s, w)
	if err != nil {
		return "", fmt.Errorf("print %s: %w", s.DataType(), err)
	}

	c := Checksum(w.String()[1:])
	fmt.Fprint(w, "*", c)

	return w.String(), nil
}
