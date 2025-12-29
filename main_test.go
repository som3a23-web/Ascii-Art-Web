package main

import (
	ascii "asciiart/features"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// helper to override the global template for deterministic responses in tests
func setTestTemplate(t *testing.T, tpl string) {
	t.Helper()
	var err error
	tmpl, err = template.New("test").Parse(tpl)
	if err != nil {
		t.Fatalf("failed to parse test template: %v", err)
	}
}

func TestHomeHandler_GET_OK(t *testing.T) {
	setTestTemplate(t, "home")
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	homeHandler(w, r)

	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}
	b, _ := io.ReadAll(res.Body)
	if string(b) != "home" {
		t.Fatalf("expected body 'home', got %q", string(b))
	}
}

func TestHomeHandler_NonRoot_404(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/not-found", nil)
	w := httptest.NewRecorder()

	homeHandler(w, r)

	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", res.StatusCode)
	}
	b, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(b), "404 Not Found") {
		t.Fatalf("expected 404 message, got %q", string(b))
	}
}

func TestHomeHandler_MethodNotAllowed_405(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	homeHandler(w, r)

	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", res.StatusCode)
	}
	b, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(b), "405 Method Not Allowed") {
		t.Fatalf("expected 405 message, got %q", string(b))
	}
}

func TestAsciiHandler_MethodNotAllowed_405(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/ascii-art", nil)
	w := httptest.NewRecorder()

	asciiHandler(w, r)

	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", res.StatusCode)
	}
}

func TestAsciiHandler_MissingParams_400(t *testing.T) {
	form := url.Values{}
	form.Set("text", "Hello")
	// banner is missing
	r := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	asciiHandler(w, r)

	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.StatusCode)
	}
}

func TestAsciiHandler_InvalidBanner_400(t *testing.T) {
	form := url.Values{}
	form.Set("text", "Hello")
	form.Set("banner", "invalid")
	r := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	asciiHandler(w, r)

	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.StatusCode)
	}
}

func TestAsciiHandler_NonASCII_400(t *testing.T) {
	form := url.Values{}
	form.Set("text", "Héllo") // contains non-ASCII 'é'
	form.Set("banner", "standard")
	r := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	asciiHandler(w, r)

	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", res.StatusCode)
	}
}

func TestAsciiHandler_TooLarge_413(t *testing.T) {
	// Build a body exceeding the 10240-byte limit
	tooLarge := strings.Repeat("a", 10241)
	form := url.Values{}
	form.Set("text", tooLarge)
	form.Set("banner", "standard")
	body := form.Encode()

	r := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	asciiHandler(w, r)

	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413, got %d", res.StatusCode)
	}
}

func TestAsciiHandler_Success_200(t *testing.T) {
	// render only the Art field to compare easily
	setTestTemplate(t, "{{.Art}}")

	input := "A"
	banner := "standard"

	// Compute expected art using the same library functions
	bannerData := ascii.ReadBanner(banner)
	if bannerData == "" {
		t.Fatalf("expected banner data to be non-empty")
	}
	bannerSlice := strings.Split(bannerData, "\n")
	art, err := ascii.DrawingInput([]string{input}, bannerSlice)
	if err != nil {
		t.Fatalf("DrawingInput returned error: %v", err)
	}

	form := url.Values{}
	form.Set("text", input)
	form.Set("banner", banner)
	r := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	asciiHandler(w, r)

	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		t.Fatalf("expected status 200, got %d with body %q", res.StatusCode, string(b))
	}
	b, _ := io.ReadAll(res.Body)
	if string(b) != art {
		t.Fatalf("unexpected art output.\nExpected:\n%q\nGot:\n%q", art, string(b))
	}
}

func TestIsValidBanner(t *testing.T) {
	cases := []struct {
		name   string
		in     string
		expect bool
	}{
		{"standard valid", "standard", true},
		{"shadow valid", "shadow", true},
		{"thinkertoy valid", "thinkertoy", true},
		{"invalid", "foobar", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isValidBanner(tc.in); got != tc.expect {
				t.Fatalf("isValidBanner(%q) = %v, want %v", tc.in, got, tc.expect)
			}
		})
	}
}

func TestIsValidASCII(t *testing.T) {
	cases := []struct {
		name   string
		in     string
		expect bool
	}{
		{"printable ASCII", "Hello, World!", true},
		{"contains newline", "Line1\nLine2", true},
		{"contains carriage return", "Line1\r\n", true},
		{"non-printable", string([]byte{0x01}), false},
		{"non-ASCII rune", "Héllo", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isValidASCII(tc.in); got != tc.expect {
				t.Fatalf("isValidASCII(%q) = %v, want %v", tc.in, got, tc.expect)
			}
		})
	}
}
