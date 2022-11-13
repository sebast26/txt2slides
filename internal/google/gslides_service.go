package google

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/slides/v1"
	"net/http"
	"strings"
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
	// create presentation
	file := drive.File{Title: fmt.Sprintf("%s - test slides", prefix)}
	copiedFile, err := s.gDrive.Files.Copy(templateFileID, &file).Do()
	if err != nil {
		return errors.Wrap(err, "failed to copy from template file")
	}
	presentation, err := s.gSlides.Presentations.Get(copiedFile.Id).Do()
	if err != nil {
		return errors.Wrap(err, "failed to create presentation file")
	}

	// chunks
	// TODO: should be split by empty lines not \n
	chunks := strings.Split(content, "\n")

	_, err = s.createEmptySlides(presentation.PresentationId, chunks)
	if err != nil {
		return errors.Wrap(err, "failed to create empty slides")
	}

	return nil
}

// create empty slides
func (s *SlidesService) createEmptySlides(presentationID string, chunks []string) ([]string, error) {
	var requests []*slides.Request
	for range chunks {
		requests = append(requests, &slides.Request{CreateSlide: &slides.CreateSlideRequest{
			ObjectId:             fmt.Sprintf("slide-%s", uuid.New().String()),
			InsertionIndex:       0,
			SlideLayoutReference: &slides.LayoutReference{PredefinedLayout: "TITLE"},
		}})
	}
	result, err := s.gSlides.Presentations.BatchUpdate(presentationID, &slides.BatchUpdatePresentationRequest{Requests: requests}).Do()
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to create empty slides")
	}
	var objectIDs []string
	for _, reply := range result.Replies {
		if reply == nil {
			continue
		}
		objectIDs = append(objectIDs, reply.CreateSlide.ObjectId)
	}
	return objectIDs, nil
}

const (
	// TODO: change it
	templateFileID = "1koVVWiMZnOoTsZ1vH2UBQ7lf3YLSU7io3e2SzbNkdEc"
)
