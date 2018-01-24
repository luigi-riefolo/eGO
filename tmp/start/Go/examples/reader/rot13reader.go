package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (r rot13Reader) Read(b []byte) (int, error) {
	i := 0
	c := make([]byte, 1)
	for {
		n, err := r.r.Read(c)
		fmt.Printf("n = %v err = %v b = %v\n", n, err, c)
		if err == io.EOF {
			return i, err
		}

		if c[0] == ' ' {
			b[i] = ' '
		} else if c[0] >= 'n' {
			b[i] = c[0] - 13
		} else {
			b[i] = c[0] + 13
		}

		i++
	}
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	_, _ = io.Copy(os.Stdout, &r)
}
