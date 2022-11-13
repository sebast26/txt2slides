package google

import "fmt"

// Presentation represents presentation.
type Presentation struct {
	ID       string
	Location string
}

// NewPresentation creates new Presentation struct.
func NewPresentation(ID string) Presentation {
	return Presentation{
		ID:       ID,
		Location: fmt.Sprintf(presentationLocationTemplate, ID),
	}
}

const (
	// presentationLocationTemplate is a template used by Google Slides to access document by ID from the browser
	presentationLocationTemplate = "https://docs.google.com/presentation/d/%s/edit"
)
