package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	handlers "github.com/media_uploader/handlers"
)

// Command line flag to specify the HTTP service address.
var addr = flag.String("addr", "localhost:8080", "http service address")

// Templates for rendering HTML pages.
var streamTemplate *template.Template
var fileSelectTemplate *template.Template

func main() {
	// Parse command line flags.
	flag.Parse()

	var err error

	// Parse HTML templates for streaming and file selection.
	streamTemplate, err = template.ParseFiles("./static/stream.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	fileSelectTemplate, err = template.ParseFiles("./static/file_select.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Start the worker pool for handling tasks.
	handlers.WorkPool.Start()

	// WebSocket endpoints for handling stream and file selection.
	http.HandleFunc("/upload_stream", handlers.StreamHandler)
	http.HandleFunc("/upload_select_file", handlers.FileSelectHandler)

	// HTTP endpoints for rendering HTML pages.
	http.HandleFunc("/stream", stream)
	http.HandleFunc("/file_select", fileSelect)

	// Start the HTTP server and log any errors.
	log.Fatal(http.ListenAndServe(*addr, nil))

	// Close the worker pool when the server is shutting down.
	handlers.WorkPool.Close()
}

// stream is an HTTP handler that renders the stream HTML page.
func stream(w http.ResponseWriter, r *http.Request) {
	// Render the stream page with the WebSocket endpoint.
	streamTemplate.Execute(w, "ws://"+r.Host+"/upload_stream")
}

// fileSelect is an HTTP handler that renders the file selection HTML page.
func fileSelect(w http.ResponseWriter, r *http.Request) {
	// Render the file selection page with the WebSocket endpoint.
	fileSelectTemplate.Execute(w, "ws://"+r.Host+"/upload_select_file")
}
