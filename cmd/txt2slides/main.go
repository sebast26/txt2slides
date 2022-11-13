package main

import (
	"github.com/sebast26/txt2slides/internal/google"
	"github.com/sebast26/txt2slides/internal/stdin"
)

func main() {
	client, err := google.NewOAuthClient("credentials.json", "token.json", false)
	if err != nil {
		panic(err)
	}
	service, err := google.NewSlidesService(client.HttpClient)
	if err != nil {
		panic(err)
	}

	buf, err := stdin.ReadStdin()
	if err != nil {
		panic(err)
	}
	presentation, err := service.CreateSlides("", string(buf))
	if err != nil {
		panic(err)
	}

	println(presentation.Location)
}
