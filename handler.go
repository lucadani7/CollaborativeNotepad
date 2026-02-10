package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// upgrader is a WebSocket upgrader that configures a custom origin checker allowing all origins.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// HandleConnections establishes a WebSocket connection for collaborative editing and manages client interactions.
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	const (
		defaultDocID = "default"
		maxEditors   = 2
	)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error WS:", err)
		return
	}
	docID := r.URL.Query().Get("doc")
	if docID == "" {
		docID = defaultDocID
	}
	doc := getOrCreateDocument(docID)
	registered := false
	defer func() {
		if registered {
			removeClient(doc, conn)
		}
		if err := conn.Close(); err != nil {
			log.Printf("WS close error for %s: %v", docID, err)
		}
	}()
	doc.mu.Lock()
	if len(doc.Clients) >= maxEditors {
		doc.mu.Unlock()
		log.Printf("Document %s full. Refusing client.", docID)
		_ = conn.WriteJSON(Message{Type: "ERROR", Content: "Maximum 2 editors allowed."})
		return
	}
	doc.Clients[conn] = true
	registered = true
	initialContent := string(doc.Content)
	doc.mu.Unlock()
	log.Printf("Client connected to: %s", docID)
	_ = conn.WriteJSON(Message{Type: "INIT", Content: initialContent})
	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("Client disconnected from %s", docID)
			return
		}
		processAndBroadcast(doc, conn, msg)
	}
}

// processAndBroadcast modifies the document based on the message and broadcasts the changes to all connected clients.
func processAndBroadcast(doc *Document, sender *websocket.Conn, msg Message) {
	doc.mu.Lock()
	defer doc.mu.Unlock()
	switch msg.Type {
	case "INSERT":
		if msg.Pos <= len(doc.Content) {
			newRunes := []rune(msg.Content)
			doc.Content = append(doc.Content[:msg.Pos], append(newRunes, doc.Content[msg.Pos:]...)...)
		}
	case "DELETE":
		if msg.Pos < len(doc.Content) {
			end := msg.Pos + msg.Len
			if end > len(doc.Content) {
				end = len(doc.Content)
			}
			doc.Content = append(doc.Content[:msg.Pos], doc.Content[end:]...)
		}
	}
	for client := range doc.Clients {
		if client != sender {
			err := client.WriteJSON(msg)
			if err != nil {
				return
			}
		}
	}
}

// removeClient removes a WebSocket client from the document's Clients map in a thread-safe manner.
func removeClient(doc *Document, ws *websocket.Conn) {
	doc.mu.Lock()
	defer doc.mu.Unlock()
	delete(doc.Clients, ws)
}

// getOrCreateDocument retrieves an existing document by ID or creates a new one if it doesn't exist in the global store.
func getOrCreateDocument(id string) *Document {
	globalStore.mu.Lock()
	defer globalStore.mu.Unlock()
	if doc, exists := globalStore.Documents[id]; exists {
		return doc
	}
	newDoc := &Document{
		ID:      id,
		Content: []rune(""),
		Clients: make(map[*websocket.Conn]bool),
	}
	globalStore.Documents[id] = newDoc
	return newDoc
}
