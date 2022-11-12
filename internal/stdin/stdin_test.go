package stdin_test

import (
	"fmt"
	"github.com/sebast26/txt2gdoc/internal/stdin"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestReadStdin(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		// given
		r, w, _ := os.Pipe()
		origStdin := os.Stdin
		os.Stdin = r
		defer func() {
			os.Stdin = origStdin
		}()

		// when
		_, _ = fmt.Fprintf(w, "")
		_ = w.Close()
		buf, err := stdin.ReadStdin()

		// then
		assert.Empty(t, buf)
		assert.NoError(t, err)
	})

	t.Run("single char", func(t *testing.T) {
		// given
		r, w, _ := os.Pipe()
		origStdin := os.Stdin
		os.Stdin = r
		defer func() {
			os.Stdin = origStdin
		}()

		// when
		_, _ = fmt.Fprintf(w, "a")
		_ = w.Close()
		buf, err := stdin.ReadStdin()

		// then
		assert.NoError(t, err)
		assert.Equal(t, []byte{'a'}, buf)
	})

	t.Run("read all chars till end of the stream", func(t *testing.T) {
		// given
		r, w, _ := os.Pipe()
		origStdin := os.Stdin
		os.Stdin = r
		defer func() {
			os.Stdin = origStdin
		}()

		// when
		_, _ = fmt.Fprintf(w, "first lin\nsecond line")
		_ = w.Close()
		buf, err := stdin.ReadStdin()

		// then
		assert.NoError(t, err)
		assert.Equal(t, []byte("first lin\nsecond line"), buf)
	})
}
