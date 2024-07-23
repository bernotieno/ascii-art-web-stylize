package utils

import (
	"log"
	"net/http"
	"text/template"
)

// ServeErrorPage serves error pages based on the provided error code.
func ServeErrorPage(w http.ResponseWriter, r *http.Request, errorCode int) {
	var templateFile string
	switch errorCode {
	case http.StatusBadRequest:
		templateFile = "templates/error400.html"
	case http.StatusNotFound:
		templateFile = "templates/error404.html"
	case http.StatusMethodNotAllowed:
		templateFile = "templates/error405.html"
	case http.StatusInternalServerError:
		templateFile = "templates/error500.html"
	default:
		templateFile = "templates/error500.html"
	}

	tmpl := template.Must(template.ParseFiles(templateFile))
	w.WriteHeader(errorCode)
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Error rendering error template: %v", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

// ServeIndex handles GET requests to the root URL.
func ServeIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "templates/index.html")
	} else {
		ServeErrorPage(w, r, http.StatusNotFound)
	}
}

// ServeAbout handles GET requests to the /about URL.
func ServeAbout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ServeErrorPage(w, r, http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "templates/about.html")
}

// GenerateASCIIArt handles POST requests to generate ASCII art.
func GenerateASCIIArt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ServeErrorPage(w, r, http.StatusMethodNotAllowed)
		return
	}

	input := r.FormValue("input")
	banner := r.FormValue("banner")

	if input == "" || banner == "" {
		ServeErrorPage(w, r, http.StatusBadRequest)
		return
	}

	content, err := ReadsFile(GetFile(banner))
	if err != nil {
		log.Printf("Error reading ASCII map: %v", err)
		ServeErrorPage(w, r, http.StatusInternalServerError)
		return
	}

	contentLines := SplitFile(string(content))

	art, err := DisplayText(input, contentLines)
	if err != nil {
		ServeErrorPage(w, r, http.StatusBadRequest)
		return
	}

	data := struct {
		Art, Input string
	}{
		Input: input,
		Art:   art,
	}

	resultTemplate := template.Must(template.ParseFiles("templates/result.html"))
	if err := resultTemplate.Execute(w, data); err != nil {
		log.Printf("Error rendering template: %v", err)
		ServeErrorPage(w, r, http.StatusInternalServerError)
	}
}

// ServeError handles GET requests to the /error URL.
func ServeError(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	switch code {
	case "400":
		ServeErrorPage(w, r, http.StatusBadRequest)
	case "404":
		ServeErrorPage(w, r, http.StatusNotFound)
	case "405":
		ServeErrorPage(w, r, http.StatusMethodNotAllowed)
	case "500":
		ServeErrorPage(w, r, http.StatusInternalServerError)
	default:
		ServeErrorPage(w, r, http.StatusInternalServerError)
	}
}
