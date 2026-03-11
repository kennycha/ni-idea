package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kennycha/ni-idea/internal/store"
)

var (
	port      string
	dataDir   string
	authToken string
)

func init() {
	flag.StringVar(&port, "port", "8080", "Server port")
	flag.StringVar(&dataDir, "data", "./data", "Data directory for notes")
}

func main() {
	flag.Parse()

	// Get auth token from environment
	authToken = os.Getenv("NI_AUTH_TOKENS")
	if authToken == "" {
		log.Println("Warning: NI_AUTH_TOKENS not set, server is unprotected!")
	}

	// Ensure data directory exists
	notesDir := filepath.Join(dataDir, "notes")
	if err := os.MkdirAll(notesDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Create type directories
	for _, noteType := range store.AllTypes {
		dir := filepath.Join(notesDir, noteType.Directory())
		os.MkdirAll(dir, 0755)
	}

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/api/ping", handlePing) // No auth for health check
	mux.HandleFunc("/api/notes", authMiddleware(handleNotes))
	mux.HandleFunc("/api/notes/", authMiddleware(handleNote))
	mux.HandleFunc("/api/search", authMiddleware(handleSearch))

	addr := ":" + port
	log.Printf("Server starting on %s (data: %s)", addr, dataDir)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if authToken == "" {
			next(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		tokens := strings.Split(authToken, ",")
		valid := false
		for _, t := range tokens {
			if strings.TrimSpace(t) == token {
				valid = true
				break
			}
		}

		if !valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func handleNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	notesDir := filepath.Join(dataDir, "notes")
	notes, err := store.ListNotes(notesDir, store.ListOptions{IncludePrivate: true})
	if err != nil {
		http.Error(w, "Failed to list notes", http.StatusInternalServerError)
		return
	}

	type NoteMeta struct {
		Path    string   `json:"path"`
		Title   string   `json:"title"`
		Type    string   `json:"type"`
		Tags    []string `json:"tags"`
		Private bool     `json:"private"`
		Created string   `json:"created"`
		Updated string   `json:"updated"`
	}

	result := make([]NoteMeta, 0, len(notes))
	for _, note := range notes {
		result = append(result, NoteMeta{
			Path:    note.Path,
			Title:   note.Meta.Title,
			Type:    string(note.Meta.Type),
			Tags:    note.Meta.Tags,
			Private: note.Meta.Private,
			Created: note.Meta.Created,
			Updated: note.Meta.Updated,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handleNote(w http.ResponseWriter, r *http.Request) {
	// Extract path from URL
	path := strings.TrimPrefix(r.URL.Path, "/api/notes/")
	if path == "" {
		http.Error(w, "Path required", http.StatusBadRequest)
		return
	}

	notesDir := filepath.Join(dataDir, "notes")

	switch r.Method {
	case http.MethodGet:
		// Get note
		fullPath, err := store.ResolvePath(notesDir, path)
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		note, err := store.ReadNote(fullPath)
		if err != nil {
			http.Error(w, "Failed to read note", http.StatusInternalServerError)
			return
		}

		content, _ := store.ReadNoteFile(fullPath)

		result := struct {
			Path    string   `json:"path"`
			Title   string   `json:"title"`
			Type    string   `json:"type"`
			Tags    []string `json:"tags"`
			Private bool     `json:"private"`
			Created string   `json:"created"`
			Updated string   `json:"updated"`
			Content string   `json:"content"`
		}{
			Path:    path,
			Title:   note.Meta.Title,
			Type:    string(note.Meta.Type),
			Tags:    note.Meta.Tags,
			Private: note.Meta.Private,
			Created: note.Meta.Created,
			Updated: note.Meta.Updated,
			Content: content,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)

	case http.MethodPut:
		// Upload/update note
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}

		var note struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		if err := json.Unmarshal(body, &note); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Use path from URL if not in body
		notePath := path
		if note.Path != "" {
			notePath = note.Path
		}

		fullPath := filepath.Join(notesDir, notePath)

		// Ensure directory exists
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			http.Error(w, "Failed to create directory", http.StatusInternalServerError)
			return
		}

		// Write file
		if err := os.WriteFile(fullPath, []byte(note.Content), 0644); err != nil {
			http.Error(w, "Failed to write file", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))

	case http.MethodDelete:
		fullPath := filepath.Join(notesDir, path)
		if err := os.Remove(fullPath); err != nil {
			if os.IsNotExist(err) {
				http.Error(w, "Not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to delete", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' required", http.StatusBadRequest)
		return
	}

	notesDir := filepath.Join(dataDir, "notes")

	// Simple search (reuse existing search logic)
	notes, err := store.ListNotes(notesDir, store.ListOptions{IncludePrivate: true})
	if err != nil {
		http.Error(w, "Failed to search", http.StatusInternalServerError)
		return
	}

	type SearchResult struct {
		Path    string   `json:"path"`
		Title   string   `json:"title"`
		Type    string   `json:"type"`
		Tags    []string `json:"tags"`
		Matches []string `json:"matches"`
	}

	var results []SearchResult
	queryLower := strings.ToLower(query)

	for _, note := range notes {
		var matches []string

		// Check title
		if strings.Contains(strings.ToLower(note.Meta.Title), queryLower) {
			matches = append(matches, "Title: "+note.Meta.Title)
		}

		// Check content
		lines := strings.Split(note.Content, "\n")
		for _, line := range lines {
			if strings.Contains(strings.ToLower(line), queryLower) {
				trimmed := strings.TrimSpace(line)
				if len(trimmed) > 100 {
					trimmed = trimmed[:100] + "..."
				}
				if trimmed != "" {
					matches = append(matches, trimmed)
				}
			}
		}

		if len(matches) > 0 {
			results = append(results, SearchResult{
				Path:    note.Path,
				Title:   note.Meta.Title,
				Type:    string(note.Meta.Type),
				Tags:    note.Meta.Tags,
				Matches: matches,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
