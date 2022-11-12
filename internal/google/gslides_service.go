package google

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/slides/v1"
	"net/http"
)

// SlidesService represents service that creates Google Slides.
type SlidesService struct {
	gSlides *slides.Service
	gDrive  *drive.Service
}

// NewSlidesService creates new SlidesService.
func NewSlidesService(oauthClient *http.Client) (*SlidesService, error) {
	ctx := context.Background()
	srv, err := slides.NewService(ctx, option.WithHTTPClient(oauthClient))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Google Slides client")
	}
	drv, err := drive.NewService(ctx, option.WithHTTPClient(oauthClient))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Google Drive client")
	}
	return &SlidesService{gSlides: srv, gDrive: drv}, nil
}

func (s *SlidesService) CreateSlides(prefix, content string) error {
	return nil
}
