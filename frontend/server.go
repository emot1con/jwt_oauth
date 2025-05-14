package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// Directory to serve static files from
	dir := "."

	// Create a file server handler
	fileServer := http.FileServer(http.Dir(dir))

	// Handle all requests by serving files
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Handle SPA routing - serve index.html for any path that doesn't correspond to a file
		path, err := filepath.Abs(r.URL.Path[1:])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// If the file doesn't exist and it's not a static asset request,
		// serve index.html to support SPA routing
		if r.URL.Path != "/" && filepath.Ext(path) == "" {
			http.ServeFile(w, r, "index.html")
			return
		}

		// Otherwise, serve the requested file
		fileServer.ServeHTTP(w, r)
	})

	// Start the server on port 3000
	fmt.Println("Frontend server running at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
