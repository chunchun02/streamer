package connection

import (
	"github.com/chunchun02/streamer/internal/work"
	"github.com/google/uuid"
	"sync"
)

type Connection struct {
	ID      uuid.UUID
	Work    *work.Work
	Channel chan *work.Work
}

type Connectors struct {
	mu          sync.Mutex
	Connections []*Connection
}

func (con *Connection) Stream() {
	con.Channel <- con.Work
}

func (connectors *Connectors) Lock() {
	connectors.mu.Lock()
}

func (connectors *Connectors) Unlock() {
	connectors.mu.Unlock()
}
