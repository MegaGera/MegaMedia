package main

import (
	"log"
	"net/http"
)

func main() {
	// Routes
	// http.HandleFunc("/upload", handlers.UploadImageHandler)
	http.Handle("/megagoal/", http.StripPrefix("/megagoal/", http.FileServer(http.Dir("./static/megagoal"))))
	http.Handle("/megagera/", http.StripPrefix("/megagera/", http.FileServer(http.Dir("./static/megagera"))))

	// Start server
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
