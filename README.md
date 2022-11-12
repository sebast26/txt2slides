# txt2slides
Created Google Slides from text input

## How to create private app?

1. Create new Google Cloud app.
2. Remove all APIs
3. Add Google Drive API and Google Slides API

## How to generate `credentials.json` file?

1. Go to the [Google Cloud Platform Console](https://console.cloud.google.com/) and from the projects list, select a project or create a new one.
2. Follow [this guide](https://support.google.com/cloud/answer/6158849) from Google to create OAuth client ID. As a scope for your application
   you may select `/auth/drive` from Google Drive API and `/auth/presentations` from Google Slides API. Add you email to the tests user group. 
   For an application type you have to select `Computer`.
3. Download and save your client credentials as `credentials.json` file.

## How to generate `token.json` file?

Please use `txt2slides-setup` program. You can invoke it with `go run github.com/sebast26/txt2slides/cmd/txt2slides-setup`