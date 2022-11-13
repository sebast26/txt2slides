package google_test

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/sebast26/txt2slides/internal/google"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestNewOAuthClient(t *testing.T) {
	t.Run("credential file not found", func(t *testing.T) {
		// when
		client, err := google.NewOAuthClient("non_existing_path", "", false)

		// then
		assert.ErrorContains(t, err, "reading client secret file failed")
		assert.Empty(t, client)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		// when
		client, err := google.NewOAuthClient("testdata/credentials_invalid.json", "", false)

		// then
		assert.ErrorContains(t, err, "parsing client secret file to config failed")
		assert.Empty(t, client)
	})

	t.Run("credentials and token found", func(t *testing.T) {
		// when
		client, err := google.NewOAuthClient("testdata/credentials.json", "testdata/token.json", false)

		// then
		assert.NoError(t, err)
		assert.NotEmpty(t, client)
	})

	t.Run("token not found, no force web", func(t *testing.T) {
		// given

		// when
		client, err := google.NewOAuthClient("testdata/credentials.json", "not_existing_token", false)

		// then
		assert.ErrorContains(t, err, "count not get token")
		assert.Empty(t, client)
	})

	t.Run("get token from web, auth server error", func(t *testing.T) {
		// given
		r, w, err := os.Pipe()
		origStdin := os.Stdin
		os.Stdin = r
		defer func() {
			os.Stdin = origStdin
		}()
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
		}))
		defer ts.Close()
		err = createCredentialsFile("/tmp/cred", ts.URL)
		assert.NoError(t, err)
		defer func() {
			_ = os.Remove("/tmp/cred")
		}()

		// when
		_, _ = fmt.Fprintf(w, "a")
		_ = w.Close()
		_, err = google.NewOAuthClient("/tmp/cred", "/tmp/token.json", true)

		// then
		assert.ErrorContains(t, err, "failed to retrieve token from web")
	})

	t.Run("get token from web", func(t *testing.T) {
		// given
		r, w, err := os.Pipe()
		origStdin := os.Stdin
		os.Stdin = r
		defer func() {
			os.Stdin = origStdin
		}()
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			_, _ = w.Write([]byte("access_token=at"))
		}))
		defer ts.Close()
		err = createCredentialsFile("/tmp/cred", ts.URL)
		assert.NoError(t, err)
		defer func() {
			_ = os.Remove("/tmp/cred")
			_ = os.Remove("/tmp/token.json")
		}()

		// when
		_, _ = fmt.Fprintf(w, "a")
		_ = w.Close()
		client, err := google.NewOAuthClient("/tmp/cred", "/tmp/token.json", true)

		// then
		assert.NoError(t, err)
		assert.NotEmpty(t, client)
		b, err := os.ReadFile("/tmp/token.json")
		assert.NoError(t, err)
		spew.Dump(b)
		var tok oauth2.Token
		err = json.Unmarshal(b, &tok)
		assert.NoError(t, err)
		assert.Equal(t, "at", tok.AccessToken)
	})
}

func createCredentialsFile(filePath, testServerURL string) error {
	cred := `{"installed": {
		"client_id": "client_id.apps.googleusercontent.com",
		"project_id": "test",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "<test_server>",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_secret": "secret",
		"redirect_uris": [
		  "http://localhost"
		]
	  }
	}`
	testCred := strings.Replace(cred, "<test_server>", testServerURL, -1)
	f, err := os.Create(filePath)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.WriteString(testCred)
	return err
}
