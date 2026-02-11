package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Message Type represents the category or nature of the message, such as "INIT", "INSERT", "DELETE", or "ERROR".
type Message struct {
	Type    string `json:"type"`
	DocID   string `json:"doc_id"`
	Pos     int    `json:"pos,omitempty"`
	Content string `json:"content,omitempty"`
	Len     int    `json:"len,omitempty"`
}

// Document represents a shared text document, enabling concurrent editing and synchronized updates across multiple clients.
type Document struct {
	ID      string
	Content []rune
	Clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

// Store represents a thread-safe container for managing shared text documents.
type Store struct {
	Documents map[string]*Document
	mu        sync.RWMutex
}

var globalStore = Store{
	Documents: make(map[string]*Document),
}
