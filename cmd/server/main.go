package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	project "github.com/andreas-bauer/rejoinderoo"
	"github.com/andreas-bauer/rejoinderoo/internal/server"
)

var html *template.Template

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	var err error
	html, err = template.ParseFS(
		project.TemplateFS,
		"web/templates/*.html",
		"web/templates/components/*.html",
	)
	if err != nil {
		log.Fatalf("Error parsing web templates: %v", err)
	}

	handlers := server.NewHandler(html)

	http.Handle("/css/output.css", http.FileServer(http.FS(project.CSS)))
	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/upload", handlers.Upload)
	http.HandleFunc("/generate", handlers.Generate)

	fmt.Printf("Server running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
