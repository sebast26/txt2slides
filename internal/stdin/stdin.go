package stdin

import (
	"bufio"
	"io"
	"os"
)

// ReadStdin returns bytes read from os.Stdin or error
func ReadStdin() ([]byte, error) {
	return io.ReadAll(bufio.NewReader(os.Stdin))
}
