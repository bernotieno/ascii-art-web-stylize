package utils_test

import (
	"ascii-art/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestPrintWord(t *testing.T) {
	contentLine, err := os.ReadFile("standard.txt")
	if err != nil {
		t.Fatalf("failed to read standard.txt: %v", err)
	}
	contentLines := utils.SplitFile(string(contentLine))
	tests := []struct {
		name     string
		word     string
		expected string
	}{
		{"Single word", "Hello", " _    _          _   _          \n| |  | |        | | | |         \n| |__| |   ___  | | | |   ___   \n|  __  |  / _ \\ | | | |  / _ \\  \n| |  | | |  __/ | | | | | (_) | \n|_|  |_|  \\___| |_| |_|  \\___/  \n                                \n                                \n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.PrintWord(test.word, contentLines)
			expected := test.expected
			if result != expected {
				t.Errorf("PrintWord(%q) =\n%s\nexpected:\n%s", test.word, result, expected)
			}
		})
	}
}

func TestIsEnglish(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"Hello, World!", true},    // English sentence
		{"123", true},              // Numbers are also considered English
		{"¡Hola, mundo!", false},   // Spanish characters
		{"こんにちは、世界！", false},       // Japanese characters
		{"Привет, мир!", false},    // Russian characters
		{"", true},                 // Empty string
		{"\n\t\r", false},          // Special characters
		{"\x7F", false},            // Non-printable ASCII character
		{"😊", false},               // Emoji
		{"Hello, World! 😀", false}, // English with emoji
		{"Hello, 世界!", false},      // English with non-English characters
		{"123 456", true},          // English with numbers and spaces
	}

	for _, test := range tests {
		if output := utils.IsEnglish(test.input); output != test.expected {
			t.Errorf("IsEnglish(%q) = %v, expected %v", test.input, output, test.expected)
		}
	}
}

func TestDisplayText(t *testing.T) {
	contentLine, err := os.ReadFile("standard.txt")
	if err != nil {
		t.Fatalf("failed to read standard.txt: %v", err)
	}
	contentLines := utils.SplitFile(string(contentLine))

	tests := []struct {
		name           string
		input          string
		expectedOutput string
		expectError    bool
	}{
		{"Single word", "Hello", " _    _          _   _          \n| |  | |        | | | |         \n| |__| |   ___  | | | |   ___   \n|  __  |  / _ \\ | | | |  / _ \\  \n| |  | | |  __/ | | | | | (_) | \n|_|  |_|  \\___| |_| |_|  \\___/  \n                                \n                                \n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := utils.DisplayText(tt.input, contentLines)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Trim trailing newlines for comparison
				output = strings.TrimRight(output, "\n")
				expected := strings.TrimRight(tt.expectedOutput, "\n")

				if output != expected {
					t.Errorf("Expected output:\n%s\n\nGot:\n%s", expected, output)
				}
			}
		})
	}
}

func TestGetFile(t *testing.T) {
	tests := []struct {
		name           string
		banner         string
		expectedOutput string
	}{
		{
			name:           "Standard banner",
			banner:         "standard",
			expectedOutput: "banners/standard.txt",
		},
		{
			name:           "Thinkertoy banner",
			banner:         "thinkertoy",
			expectedOutput: "banners/thinkertoy.txt",
		},
		{
			name:           "Shadow banner",
			banner:         "shadow",
			expectedOutput: "banners/shadow.txt",
		},
		{
			name:           "Invalid banner",
			banner:         "invalid",
			expectedOutput: "Invalid bannerfile name",
		},
		{
			name:           "Empty input",
			banner:         "",
			expectedOutput: "Invalid bannerfile name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.GetFile(tt.banner)
			if result != tt.expectedOutput {
				t.Errorf("GetFile(%q) = %q, want %q", tt.banner, result, tt.expectedOutput)
			}
		})
	}
}

func TestServeIndex(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"Valid GET /", http.MethodGet, "/", http.StatusOK},
		{"Invalid Method", http.MethodPost, "/", http.StatusBadRequest},
		{"Not Found", http.MethodGet, "/invalid", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(utils.ServeIndex)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestGenerateASCIIArt(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		input          string
		banner         string
		expectedStatus int
		expectedBody   string
	}{
		{"Invalid Method", http.MethodGet, "", "", http.StatusMethodNotAllowed, ""},
		{"Empty Input", http.MethodPost, "", "standard", http.StatusBadRequest, "Input is required"},
		{"Invalid Banner", http.MethodPost, "Hello", "invalid", http.StatusInternalServerError, "Error reading ASCII map"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := strings.NewReader("input=" + tt.input + "&banner=" + tt.banner)
			req, err := http.NewRequest(tt.method, "/generate", form)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(utils.GenerateASCIIArt)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if tt.expectedBody != "" && !strings.Contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
