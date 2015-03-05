package gist6096872

import (
	"bufio"
	"io"
)

type ChanWriter chan []byte

func (cw ChanWriter) Write(p []byte) (n int, err error) {
	// TODO: Copy the slice contents rather than sending the original, as it may get modified?
	cw <- p
	return len(p), nil
}

// Credit to Tarmigan.
// TODO: Delete this? It's not used anywhere.
func ByteReader(r io.Reader) <-chan []byte {
	ch := make(chan []byte)
	go func() {
		for {
			buf := make([]byte, 2048)
			s := 0

		Inner:
			for {
				n, err := r.Read(buf[s:])
				if n > 0 {
					ch <- buf[s : s+n]
					s += n
				}
				if err != nil {
					close(ch)
					return
				}
				if s >= len(buf) {
					break Inner
				}
			}
		}
	}()

	return ch
}

func LineReader(r io.Reader) <-chan []byte {
	ch := make(chan []byte)
	go func() {
		br := bufio.NewReader(r)

		for {
			line, err := br.ReadBytes('\n')
			if err != nil {
				ch <- line
				close(ch)
				return
			}

			ch <- line[:len(line)-1] // Trim last newline.
		}
	}()

	return ch
}
