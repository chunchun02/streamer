package connection

import (
	"github.com/chunchun02/streamer/internal/work"
	"github.com/google/uuid"
	"sync"
)

// Connection presents for a connection between user and server
type Connection struct {
	ID      uuid.UUID       // Unique UUID for each connection
	Channel chan *work.Work // Channel to send data to user
}

// Connectors present for a list of all current connecting users
type Connectors struct {
	mu          sync.Mutex    // Mutex to lock/unlock while working as multi-threading
	Connections []*Connection // List of all current connections
}

// Lock locks the Connectors
func (connectors *Connectors) Lock() {
	connectors.mu.Lock()
}

// Unlock unlocks the Connectors
func (connectors *Connectors) Unlock() {
	connectors.mu.Unlock()
}
