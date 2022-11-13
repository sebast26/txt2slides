package google_test

import (
	"github.com/sebast26/txt2slides/internal/google"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPresentation(t *testing.T) {
	t.Run("location", func(t *testing.T) {
		// when
		presentation := google.NewPresentation("some-id")

		// then
		assert.Equal(t, "some-id", presentation.ID)
		assert.Equal(t, "https://docs.google.com/presentation/d/some-id/edit", presentation.Location)
	})
}
