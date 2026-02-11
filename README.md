# Collaborative Notepad
A high-performance client-server application that allows simultaneous text editing by multiple users. Built with **Go (Golang)** for the backend and **Vanilla JavaScript** for the frontend, containerized with **Docker**. This project implements a distributed system where multiple clients can edit the same text document in real-time. It handles concurrency, synchronization, and conflict resolution without relying on heavy external frameworks.

---

## Disclaimer
* This application is intended for **educational purposes only** to demonstrate distributed system concepts (WebSockets, Concurrency in Go, Docker).
* It is **not** a production-ready replacement for professional tools like Google Docs, VS Code Live Share, or similar solutions. It implements a basic synchronization logic suitable for learning, but lacks advanced features found in commercial software (e.g., complex Operational Transformation, OT/CRDT algorithms, persistent database storage, or user authentication).

---

## Key Features
* **Real-Time Synchronization:** Uses WebSockets for low-latency bi-directional communication.
* **Concurrency Control:** Strictly limits access to **maximum 2 concurrent editors** per document.
* **Conflict Resolution:** Implements a diff-based synchronization algorithm. The client calculates changes (insertions/deletions) and the server broadcasts them, allowing clients to adjust their cursors automatically.
* **Multi-Document Support:** Supports multiple separate editing sessions via URL parameters (e.g., `?doc=file1`).
* **Thread Safety:** The Go server uses `sync.RWMutex` to ensure memory safety during concurrent writes.

---

## Installation & Getting Started
You can run this application using **Docker** (recommended) or by compiling it manually.

### Option A: using Docker
This ensures the application runs in an isolated environment compatible with any OS (Windows, Mac, Linux).

1. Clone the repository:
   ```bash
   git clone https://github.com/lucadani7/CollaborativeNotepad
   cd CollaborativeNotepad
   ```
2. Run the container:
   ```bash
   docker compose up --build
   ```
3. The server will start at: `http://localhost:8080`

### Option B: Manual Compilation
If you have Go installed locally (version 1.25+ recommended).

1. Clone the repository:
   ```bash
   git clone https://github.com/lucadani7/CollaborativeNotepad
   cd CollaborativeNotepad
   ```
2. Initialize dependencies:
   ```bash
   go mod download
   ```
3. Build the executable:
   ```bash
   go build -o server .
   ```
4. Run the application:
   ```bash
   ./server
   ```

---

## Architecture Details

### The Protocol
Communication happens via JSON messages over WebSockets.
* INSERT: Sent when a user adds text. Contains position and content.
* DELETE: Sent when a user removes text. Contains position and length.
* INIT: Sent by the server to a new client with the full document state.

### Concurrency Model
* The server maintains a map of Document structs.
* Each document is protected by a `sync.Mutex` to prevent race conditions.
* When a client sends an edit, the server applies it to its in-memory "Source of Truth" and broadcasts the delta (the change) to the other connected client.
---

## License
This project is licensed under the Apache-2.0 License.
