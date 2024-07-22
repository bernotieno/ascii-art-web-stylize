package utils

import (
	"log"
	"net/http"
	"text/template"
)

// ServeIndex handles GET requests to the root URL
func ServeIndex(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	// Serve the index HTML file if the URL path is "/"
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "templates/index.html")
	} else {
		// Return a 404 error for any other path
		http.NotFound(w, r)
	}
}

// GenerateASCIIArt handles POST requests to generate ASCII art
func GenerateASCIIArt(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the input text and banner type from the form data
	input := r.FormValue("input")
	banner := r.FormValue("banner")

	// Validate that input text is provided
	if input == "" || banner == "" {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	// Read the ASCII map file based on the banner type
	content, err := ReadsFile(GetFile(banner))
	if err != nil {
		log.Printf("Error reading ASCII map: %v", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Split the content of the file into lines
	contentLines := SplitFile(string(content))

	// Generate the ASCII art from the input text and the content lines
	art, err := DisplayText(input, contentLines)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	// Define the data to be passed to the template
	data := struct {
		Art, Input string
	}{
		Input: input,
		Art: art,
	}

	// Parse the result HTML template at initialization
	resultTemplate := template.Must(template.ParseFiles("templates/result.html"))

	// Render the result template with the generated ASCII art
	if err := resultTemplate.Execute(w, data); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
}
