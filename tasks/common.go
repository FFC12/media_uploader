package tasks

import (
	"encoding/json"
	"sync"
)

// FirstChunk represents a data structure for a task's first chunk that comes from the client.
type FirstChunk struct {
	Video    bool   `json:"video"`
	MimeType string `json:"mimeType"`
	MediaId  string `json:"mediaId"`
}

// Serializer is a goroutine-safe struct that facilitates the concurrent serialization and deserialization of JSON data.
type Serializer struct {
	mu sync.Mutex
}

// Serialize serializes a FirstChunk instance to JSON format in a goroutine-safe manner.
func (s *Serializer) Serialize(firstChunk FirstChunk) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return json.Marshal(firstChunk)
}

// Deserialize deserializes JSON data into a FirstChunk instance in a goroutine-safe manner.
func (s *Serializer) Deserialize(data []byte) (FirstChunk, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var firstChunk FirstChunk
	err := json.Unmarshal(data, &firstChunk)
	return firstChunk, err
}

// JsonSerializer is a shared instance of Serializer that can be used safely across multiple goroutines.
var JsonSerializer = &Serializer{}
