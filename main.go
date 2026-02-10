package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", HandleConnections)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	log.Println("------------------------------------------------")
	log.Println("Server is running on: http://localhost:8080")
	log.Println("Press Ctrl+C to disconnect.")
	log.Println("------------------------------------------------")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server error: ", err)
	}
}
