package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string `json:"type"`
	DocID   string `json:"doc_id"`
	Pos     int    `json:"pos,omitempty"`
	Content string `json:"content,omitempty"`
	Len     int    `json:"len,omitempty"`
}

type Document struct {
	ID      string
	Content []rune
	Clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

type Store struct {
	Documents map[string]*Document
	mu        sync.RWMutex
}

var globalStore = Store{
	Documents: make(map[string]*Document),
}
