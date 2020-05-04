package connection

import (
	"github.com/chunchun02/streamer/internal/work"
	"github.com/google/uuid"
	"sync"
)

type Connection struct {
	ID      uuid.UUID
	Channel chan *work.Work
}

type Connectors struct {
	mu          sync.Mutex
	Connections []*Connection
}

func (connectors *Connectors) Lock() {
	connectors.mu.Lock()
}

func (connectors *Connectors) Unlock() {
	connectors.mu.Unlock()
}
