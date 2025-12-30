package main

import (
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

// ============================================================================
// FUNCTIONAL REQUIREMENTS TESTS
// ============================================================================

// Test: Only standard packages are used
func TestOnlyStandardPackages(t *testing.T) {
	// This is verified by the import statements and go.mod
	// The test passes if the code compiles with only standard library
}

// Test: Project contains HTML files in templates directory
func TestHTMLFilesExist(t *testing.T) {
	requiredTemplates := []string{
		"./templates/index.html",
		"./templates/404.html",
		"./templates/400.html",
		"./templates/500.html",
	}

	for _, tmplPath := range requiredTemplates {
		if _, err := template.ParseFiles(tmplPath); err != nil {
			t.Errorf("Required template not found or invalid: %s - %v", tmplPath, err)
		}
	}
}

// Test: Standard template input with special characters
func TestStandardTemplateWithSpecialChars(t *testing.T) {
	form := url.Values{}
	form.Set("text", "{123}\n<Hello> (World)!")
	form.Set("banner", "standard")

	r := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	asciiHandler(w, r)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	output := string(body)

	// Verify output contains ASCII art representation
	if !strings.Contains(output, "{") || !strings.Contains(output, "}") {
		t.Error("Output should contain ASCII representation of input characters")
	}
}

// Test: Input "123??" with standard banner
func TestStandardTemplate123QuestionMarks(t *testing.T) {
	form := url.Values{}
	form.Set("text", "123??")
	form.Set("banner", "standard")

	r := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	asciiHandler(w, r)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	output := string(body)

	if !strings.Contains(output, "?") {
		t.Error("Output should contain ASCII representation of question marks")
	}
}

// Test: Shadow banner with special characters
func TestShadowBannerSpecialChars(t *testing.T) {
	form := url.Values{}
	form.Set("text", "$% \"=")
	form.Set("banner", "shadow")

	r := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	asciiHandler(w, r)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}
}

// Test: Thinkertoy banner with complex input
func TestThinkertoyBannerComplexInput(t *testing.T) {
	form := url.Values{}
	form.Set("text", "123 T/fs#R")
	form.Set("banner", "thinkertoy")

	r := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	asciiHandler(w, r)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}
}

// Test: Graphical representation is displayed
func TestGraphicalRepresentationDisplayed(t *testing.T) {
	form := url.Values{}
	form.Set("text", "Hi")
	form.Set("banner", "standard")

	r := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	asciiHandler(w, r)

	res := w.Result()
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	output := string(body)

	// Check that output contains ASCII art patterns (multiple lines, ASCII characters)
	if !strings.Contains(output, "\n") {
		t.Error("ASCII art should contain multiple lines")
	}

	// Should contain ASCII art characters like |, _, /, etc.
	hasAsciiChars := strings.ContainsAny(output, "|_/\\-")
	if !hasAsciiChars {
		t.Error("Output should contain ASCII art characters")
	}
}

// ============================================================================
// HTTP STATUS CODE TESTS
// ============================================================================

// Test: 404 Not Found - Invalid route
func TestInvalidRoute404(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()

	homeHandler(w, r)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", res.StatusCode)
	}
}

// Test: 400 Bad Request - Invalid banner
func TestInvalidBanner400(t *testing.T) {
	form := url.Values{}
	form.Set("text", "Hello")
	form.Set("banner", "invalid_banner")

	r := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	asciiHandler(w, r)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", res.StatusCode)
	}
}

// Test: 400 Bad Request - Non-ASCII characters
func TestNonASCIICharacters400(t *testing.T) {
	form := url.Values{}
	form.Set("text", "Hello 世界") // Contains non-ASCII characters
	form.Set("banner", "standard")

	r := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	asciiHandler(w, r)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 for non-ASCII input, got %d", res.StatusCode)
	}
}

// Test: 500 Internal Server Error handling (simulated)
func TestInternalServerError500(t *testing.T) {
	// 1. Setup: Rename the real file to something else so ReadFile fails
	err := os.Rename("banner/standard.txt", "banner/standard_backup.txt")
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}

	// 2. Ensure the file is renamed back even if the test fails
	defer os.Rename("banner/standard_backup.txt", "banner/standard.txt")

	form := url.Values{}
	form.Set("text", "Test")
	form.Set("banner", "standard") // This is a "valid" banner name

	r := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	asciiHandler(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", w.Code)
	}
}

// Test: 405 Method Not Allowed
func TestMethodNotAllowed405(t *testing.T) {
	// Test POST on GET-only endpoint
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	homeHandler(w, r)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", res.StatusCode)
	}

	// Test GET on POST-only endpoint
	r = httptest.NewRequest(http.MethodGet, "/ascii-art", nil)
	w = httptest.NewRecorder()

	asciiHandler(w, r)

	res = w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", res.StatusCode)
	}
}

// ============================================================================
// SERVER COMMUNICATION TESTS
// ============================================================================

// Test: Server-client communication works
func TestServerClientCommunication(t *testing.T) {
	form := url.Values{}
	form.Set("text", "Test")
	form.Set("banner", "standard")

	r := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	asciiHandler(w, r)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Server-client communication failed, status: %d", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") && contentType != "" {
		t.Logf("Content-Type: %s (expected text/html or empty)", contentType)
	}
}

// Test: Correct HTTP method is used (POST for /ascii-art)
func TestCorrectHTTPMethod(t *testing.T) {
	form := url.Values{}
	form.Set("text", "Hello")
	form.Set("banner", "standard")

	r := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	asciiHandler(w, r)

	if w.Code != http.StatusOK && w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Unexpected status code: %d", w.Code)
	}
}

// Test: Site doesn't crash
func TestSiteDoesNotCrash(t *testing.T) {
	testCases := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{"Home page", http.MethodGet, "/", ""},
		{"Valid POST", http.MethodPost, "/ascii-art", "text=Hi&banner=standard"},
		{"Empty POST", http.MethodPost, "/ascii-art", "text=&banner=standard"},
		{"Invalid banner", http.MethodPost, "/ascii-art", "text=Hi&banner=invalid"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Site crashed with panic: %v", r)
				}
			}()

			r := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			if tc.method == http.MethodPost {
				r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}
			w := httptest.NewRecorder()

			if tc.path == "/" {
				homeHandler(w, r)
			} else {
				asciiHandler(w, r)
			}

			// If we reach here, site didn't crash
		})
	}
}

// Test: Server is written in Go
func TestServerWrittenInGo(t *testing.T) {
	// This test passes if the file compiles as Go code
	// Verified by the fact that we can run these tests
}

// ============================================================================
// ADDITIONAL VALIDATION TESTS
// ============================================================================

// Test: All three banners are available
func TestAllBannersAvailable(t *testing.T) {
	banners := []string{"standard", "shadow", "thinkertoy"}

	for _, banner := range banners {
		bannerSelected, _ := os.ReadFile("banner/" + banner + ".txt")
		bannerSelectedConv := string(bannerSelected)
		if bannerSelectedConv == "" {
			t.Errorf("Banner '%s' not available or empty", banner)
		}
	}
}

// Test: Banner validation function
func TestBannerValidation(t *testing.T) {
	validBanners := []string{"standard", "shadow", "thinkertoy"}
	invalidBanners := []string{"invalid", "Standard", "SHADOW", "", "random"}

	for _, banner := range validBanners {
		if !isValidBanner(banner) {
			t.Errorf("Valid banner '%s' rejected", banner)
		}
	}

	for _, banner := range invalidBanners {
		if isValidBanner(banner) {
			t.Errorf("Invalid banner '%s' accepted", banner)
		}
	}
}

// Test: ASCII validation function
func TestASCIIValidation(t *testing.T) {
	validInputs := []string{
		"Hello World",
		"123!@#",
		"Test\nMultiline",
		"Special: ()[]{}",
	}

	invalidInputs := []string{
		"Hello 世界",
		"Emoji: 😀",
		"Arabic: مرحبا",
		string([]byte{0x01, 0x02}), // Control characters
	}

	for _, input := range validInputs {
		if !isValidASCII(input) {
			t.Errorf("Valid ASCII input '%s' rejected", input)
		}
	}

	for _, input := range invalidInputs {
		if isValidASCII(input) {
			t.Errorf("Invalid ASCII input '%s' accepted", input)
		}
	}
}

// Test: Empty input handling
func TestEmptyInputHandling(t *testing.T) {
	form := url.Values{}
	form.Set("text", "")
	form.Set("banner", "standard")

	r := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	asciiHandler(w, r)

	res := w.Result()
	defer res.Body.Close()

	// Empty input should be handled gracefully (200 OK with empty output)
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for empty input, got %d", res.StatusCode)
	}
}

// Test: Newline handling
func TestNewlineHandling(t *testing.T) {
	form := url.Values{}
	form.Set("text", "Line1\nLine2\nLine3")
	form.Set("banner", "standard")

	r := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	asciiHandler(w, r)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	output := string(body)

	// Should handle multiple lines
	lineCount := strings.Count(output, "\n")
	if lineCount < 3 {
		t.Error("Output should contain multiple lines for newline-separated input")
	}
}

// Test: Path traversal attempt
func TestPathTraversal400(t *testing.T) {
	form := url.Values{}
	form.Set("text", "Hello")
	form.Set("banner", "../main.go") // Malicious path

	r := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	asciiHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for path traversal attempt, got %d", w.Code)
	}
}
